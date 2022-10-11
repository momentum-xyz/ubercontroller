-- Core Plugin
INSERT INTO plugin(plugin_id, plugin_name, description, created_at, updated_at)
VALUES ('F0F0F0F0-0F0F-4FF0-AF0F-F0F0F0F0F0F0', 'core', 'backend core', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (plugin_id) DO UPDATE SET plugin_name = excluded.plugin_name,
                                      description = excluded.description;

INSERT INTO attribute
    (plugin_id, attribute_name, description)
VALUES ('F0F0F0F0-0F0F-4FF0-AF0F-F0F0F0F0F0F0', 'node_settings', '')
ON CONFLICT (plugin_id,attribute_name)
    DO UPDATE SET description = excluded.description;

INSERT INTO attribute
    (plugin_id, attribute_name, description, options)
VALUES ('F0F0F0F0-0F0F-4FF0-AF0F-F0F0F0F0F0F0', 'name', 'Space name', '{"render_type":"texture"}')
ON CONFLICT (plugin_id,attribute_name)
    DO UPDATE SET description = excluded.description,
                  options='{"render_type":"texture"}';

INSERT INTO attribute
(plugin_id, attribute_name, description, options)
VALUES ('F0F0F0F0-0F0F-4FF0-AF0F-F0F0F0F0F0F0', 'description', 'Space name', '{"render_type":"texture"}')
ON CONFLICT (plugin_id,attribute_name)
    DO UPDATE SET description = excluded.description,
                  options='{"render_type":"texture"}';

-- Dashboard Plugin
INSERT INTO plugin (plugin_id, plugin_name, description, created_at, updated_at)
VALUES ('220578C8-5FEC-42C8-ADE1-14D970E714BD', 'dashboard', 'dashboard plugin', CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP) ON CONFLICT (plugin_id) DO
    UPDATE
    SET plugin_name = excluded.plugin_name, description = excluded.description;

INSERT INTO attribute
(plugin_id, attribute_name,description, options)
VALUES
    ('220578C8-5FEC-42C8-ADE1-14D970E714BD', 'third', 'Third screen for space','{"render_type":"texture"}')
ON CONFLICT (plugin_id,attribute_name)
    DO UPDATE SET
                  description = excluded.description,
                  options='{"render_type":"texture","content_type":"image"}';

INSERT INTO attribute
(plugin_id, attribute_name,description, options)
VALUES
    ('220578C8-5FEC-42C8-ADE1-14D970E714BD', 'poster', 'Poster for space','{"render_type":"texture"}')
ON CONFLICT (plugin_id,attribute_name)
    DO UPDATE SET
                  description = excluded.description,
                  options='{"render_type":"texture","content_type":"image"}';

INSERT INTO attribute
(plugin_id, attribute_name,description, options)
VALUES
    ('220578C8-5FEC-42C8-ADE1-14D970E714BD', 'meme', 'Meme for space','{"render_type":"texture"}')
ON CONFLICT (plugin_id,attribute_name)
    DO UPDATE SET
                  description = excluded.description,
                  options='{"render_type":"texture","content_type":"image"}';

INSERT INTO attribute
(plugin_id, attribute_name,description, options)
VALUES
    ('220578C8-5FEC-42C8-ADE1-14D970E714BD', 'video', 'Video for space','{"render_type":"texture","content_type":"video"}')
ON CONFLICT (plugin_id,attribute_name)
    DO UPDATE SET
                  description = excluded.description,
                  options='{"render_type":"texture","content_type":"video"}';

INSERT INTO attribute
(plugin_id, attribute_name,description, options)
VALUES
    ('220578C8-5FEC-42C8-ADE1-14D970E714BD', 'problem', 'Problem for space','{"render_type":"texture","content_type":"text"}')
ON CONFLICT (plugin_id,attribute_name)
    DO UPDATE SET
                  description = excluded.description,
                  options='{"render_type":"texture","content_type":"text"}';

INSERT INTO attribute
(plugin_id, attribute_name,description, options)
VALUES
    ('220578C8-5FEC-42C8-ADE1-14D970E714BD', 'solution', 'solution for space',' {"render_type":"texture","content_type":"text"}')
ON CONFLICT (plugin_id,attribute_name)
    DO UPDATE SET
                  description = excluded.description,
                  options='{"render_type":"texture","content_type":"text"}';


INSERT INTO attribute
(plugin_id, attribute_name,description, options)
VALUES
    ('220578C8-5FEC-42C8-ADE1-14D970E714BD', 'tile', 'tile for space','{"render_type":"texture"}')
ON CONFLICT (plugin_id,attribute_name)
    DO UPDATE SET
    description = excluded.description;



-- High Five plugin
INSERT INTO plugin (plugin_id, plugin_name, description, created_at, updated_at)
VALUES ('3253D616-215F-47A9-BA9D-93185EB3E6B5', 'high-five', 'High-Five plugin', CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP) ON CONFLICT (plugin_id) DO
    UPDATE
    SET plugin_name = excluded.plugin_name, description = excluded.description;

INSERT INTO attribute
(plugin_id, attribute_name,description)
VALUES
    ('3253D616-215F-47A9-BA9D-93185EB3E6B5', 'count', 'High5s given')
ON CONFLICT (plugin_id,attribute_name)
    DO UPDATE SET
    description = excluded.description;

-- Miro Plugin
INSERT INTO plugin (plugin_id, plugin_name, description, created_at, updated_at)
VALUES ('24071066-E8C6-4692-95B5-AE2DC3ED075C', 'miro', 'Miro dashboard integration plugin', CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP) ON CONFLICT (plugin_id) DO
    UPDATE
    SET plugin_name = excluded.plugin_name, description = excluded.description;
--

INSERT INTO attribute
(plugin_id, attribute_name,description)
VALUES
    ('24071066-E8C6-4692-95B5-AE2DC3ED075C', 'config', 'Miro configuration')
ON CONFLICT (plugin_id,attribute_name)
    DO UPDATE SET
    description = excluded.description;

INSERT INTO attribute
(plugin_id, attribute_name,description)
VALUES
    ('24071066-E8C6-4692-95B5-AE2DC3ED075C', 'board', 'Miro board')
ON CONFLICT (plugin_id,attribute_name)
    DO UPDATE SET
    description = excluded.description;


-- Kusama wallet plugin
INSERT INTO plugin (plugin_id, plugin_name, description, created_at, updated_at)
VALUES ('86DC3AE7-9F3D-42CB-85A3-A71ABC3C3CB8', 'kusama_wallet', 'Kusama/Substrate Wallet Plugin', CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP) ON CONFLICT (plugin_id) DO
    UPDATE
    SET plugin_name = excluded.plugin_name, description = excluded.description;

INSERT INTO attribute
(plugin_id, attribute_name,description)
VALUES
    ('86DC3AE7-9F3D-42CB-85A3-A71ABC3C3CB8', 'wallet', 'Kusama/Substrate wallet')
ON CONFLICT (plugin_id,attribute_name)
    DO UPDATE SET
    description = excluded.description;
