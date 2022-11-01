update plugin set meta = '{
  "scriptUrl": "https://dev2.odyssey.ninja/plugins/miro/remoteEntry.js",
  "scopeName": "plugin_miro",
  "name": "Miro"
}' where plugin_id='24071066-E8C6-4692-95B5-AE2DC3ED075C';

update plugin set options = '{
  "subPath": "miro",
  "exact": true,
  "iconName": "miro"
}' where plugin_id='24071066-E8C6-4692-95B5-AE2DC3ED075C';

drop index idx_20929_index_attributes_name;