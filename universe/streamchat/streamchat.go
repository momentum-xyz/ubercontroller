// Text chat functionality backed by getstream.io chat service.
//
// TODO: convert to plugin.
package streamchat

import (
	"time"

	stream "github.com/GetStream/stream-chat-go/v6"
	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

// Currently manage a single channel type
const MomentumChannelType = "momentum-world"

// String we pass into streamchat as the top admin user ID.
const systemUser = "00000000-0000-0000-0000-000000000000"

type StreamChat struct {
	node      universe.Node
	ctx       context.Context
	log       *zap.SugaredLogger
	apiKey    string
	apiSecret string
	client    *stream.Client
}

func NewStreamChat() *StreamChat {
	return &StreamChat{}
}

func (s *StreamChat) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.LoggerContextKey))
	}
	cfg := utils.GetFromAny(ctx.Value(types.ConfigContextKey), (*config.Config)(nil))
	if cfg == nil {
		return errors.Errorf("failed to get config from context: %T", ctx.Value(types.ConfigContextKey))
	}
	apiKey := cfg.Streamchat.APIKey
	apiSecret := cfg.Streamchat.APISecret

	s.node = universe.GetNode()
	s.ctx = ctx
	s.log = log
	s.apiKey = apiKey
	s.apiSecret = apiSecret
	return nil
}

func (s *StreamChat) Load() error {
	s.log.Debugf("Loading streamchat service...")
	if s.apiKey == "" || s.apiSecret == "" {
		// TODO: for now ignoring, should be a plugin anyway (and need to distribute keys to all devs)
		// Should become runtime configuration a node/world admin can configure.
		s.log.Warn("No streamchat API key and secret, not loading.")
		return nil
	}

	client, err := stream.NewClient(s.apiKey, s.apiSecret)
	if err != nil {
		return err
	}
	s.client = client

	if err := s.updateSettings(); err != nil {
		return err
	}

	universe.GetNode().AddAPIRegister(s)

	return nil
}

// Update the remote service settings.
//
// We manage the settings here, instead of manually keeping this in sync through the website.
func (s *StreamChat) updateSettings() error {
	settings := stream.NewAppSettings()
	// .app scope permissions, reset default perms for 'external' users.
	// (leaves admin/global_admin and moderator roles with the defaults).
	settings.SetGrants(
		map[string][]string{
			"anonymous":         {},
			"guest":             {},
			"channel_member":    {},
			"channel_moderator": {},
			"user":              {"update-user-owner"},
		},
	)
	_, err := s.client.UpdateAppSettings(s.ctx, settings)
	if err != nil {
		return err
	}
	if err := s.ensureChannelType(MomentumChannelType); err != nil {
		return err
	}

	commonGrants := []string{"read-channel", "create-message", "add-links"}
	channelConfig := map[string]interface{}{
		"commands":           []string{},
		"max_message_length": 280,
		"quotes":             false,
		"replies":            false,
		"search":             false,
		"typing_events":      false,
		"uploads":            false,
		"url_enrichment":     false,
		"reactions":          false,
		"mutes":              false,
		"push_notifications": false,
		"grants": map[string][]string{
			"guest":             {},
			"anonymous":         {},
			"user":              {},
			"channel_member":    commonGrants,
			"channel_moderator": append(commonGrants, "delete-message"),
		},
	}
	_, err = s.client.UpdateChannelType(s.ctx, MomentumChannelType, channelConfig)
	if err != nil {
		return err
	}

	return nil
}

// Check if channel type exists, create it if not.
func (s *StreamChat) ensureChannelType(name string) error {
	_, err := s.client.GetChannelType(s.ctx, name)
	if err != nil {
		// Gonna assume it does not exist.
		if err := s.createChannelType(MomentumChannelType); err != nil {
			return err
		}
	}
	return nil

}

func (s *StreamChat) createChannelType(name string) error {
	s.log.Debugf("Creating streamchat channel %s", MomentumChannelType)
	chType := &stream.ChannelType{
		ChannelConfig: stream.DefaultChannelConfig,
	}
	chType.ChannelConfig.Name = MomentumChannelType
	_, err := s.client.CreateChannelType(s.ctx, chType)
	if err != nil {
		return err
	}
	return nil
}

// Get StreamChat authentication token.
// Creates or update the user at the streamchat side.
func (s *StreamChat) GetToken(ctx context.Context, user universe.User) (string, error) {
	userUUID := user.GetID().String()
	profile := user.GetProfile()
	var name string
	if profile.Name != nil {
		name = *profile.Name
	}
	streamUser := &stream.User{
		ID:   userUUID,
		Name: name,
	}
	_, err := s.client.UpsertUser(ctx, streamUser)
	if err != nil {
		return "", err
	}
	issuedAt := time.Now().UTC()
	expiredAt := issuedAt.Add(24 * time.Hour) // TODO: some sane value
	token, err := s.client.CreateToken(userUUID, expiredAt, issuedAt)
	if err != nil {
		return "", err
	}
	return token, nil
}

// Get or create a channel for given space.
func (s *StreamChat) GetChannel(ctx context.Context, space universe.Space) (*stream.Channel, error) {

	chanID := space.GetID().String()
	userID := systemUser
	response, err := s.client.CreateChannel(ctx, MomentumChannelType, chanID, userID, nil)
	if err != nil {
		return nil, err
	}
	return response.Channel, nil

}

// Add user as a member of a channel.
func (s *StreamChat) MakeMember(ctx context.Context, channel *stream.Channel, user universe.User) error {
	userUUID := user.GetID().String()
	_, err := channel.AddMembers(
		ctx,
		[]string{userUUID},
		stream.AddMembersWithRolesAssignment(
			[]*stream.RoleAssignment{{
				ChannelRole: "channel_member",
				UserID:      userUUID,
			}},
		),
	)
	if err != nil {
		return err
	}
	return nil
}

// Remove user as a member from a streamchat channel.
func (s *StreamChat) RemoveMember(ctx context.Context, channel *stream.Channel, user universe.User) error {
	userUUID := user.GetID().String()
	_, err := channel.RemoveMembers(ctx, []string{userUUID}, nil)
	if err != nil {
		return err
	}
	return nil
}
