package node

import (
	"context"
	"fmt"

	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

const query = `
UPDATE node_attribute
SET
  value = (CASE
    WHEN value->$1 IS NOT NULL
      THEN jsonb_set(value,
        $2,
        (COALESCE(value#>$2, '0')::int + 1) ::text::jsonb
        )
    WHEN value->$1 IS NULL
      THEN jsonb_insert(value, $3, $4)
    END)
WHERE plugin_id = $5 AND attribute_name = $6
;
`

// Usage tracking of AI generation calls.
func (n *Node) trackAIUsage(ctx context.Context, provider string, userID umid.UMID) {
	// Quick hack...
	// Yes, i'm faster at hacky SQL then go+attrs modify functions.
	// Note: so this attr is not kept up to date in-mem!
	n.log.Debugf("AI tracker: %s %s", provider, userID)
	conn := n.db.GetCommonDB().GetConnection()
	jsonPath := fmt.Sprintf("{\"%s\",\"%s\"}", userID, provider)
	jsonNestedPath := fmt.Sprintf("{\"%s\"}", userID)
	jsonInitVal := fmt.Sprintf("{\"%s\": 1}", provider)
	if _, err := conn.Exec(
		ctx, query,
		userID.String(), jsonPath, jsonNestedPath, jsonInitVal,
		universe.GetSystemPluginID(), "tracker_ai_usage",
	); err != nil {
		//ignore tracking errors! This is our loss, not the users :)
		n.log.Error(err)
	}
}
