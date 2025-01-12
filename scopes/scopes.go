package scopes

import "encoding/json"

type Scope struct {
	Value     string
	ForServer bool
} // @name Scopes

var (
	ScopeAdmin       = registerNonServerScope("admin")
	ScopeLogin       = registerNonServerScope("login")        //can you log in
	ScopeOAuth2Auth  = registerNonServerScope("oauth2.auth")  //can you validate user credentials over OAuth2
	ScopeNodesView   = registerNonServerScope("nodes.view")   //can you globally view nodes
	ScopeNodesCreate = registerNonServerScope("nodes.create") //can you create nodes
	ScopeNodesEdit   = registerNonServerScope("nodes.edit")   //can you edit an existing node
	ScopeNodesDelete = registerNonServerScope("nodes.delete") //can you delete a node
	ScopeNodesDeploy = registerNonServerScope("nodes.deploy") //can you deploy the node (this has secret info, which is why it's special)
	ScopeSelfEdit    = registerNonServerScope("self.edit")    //can you manage your own account
	ScopeSelfClients = registerNonServerScope("self.clients") //can the user create and manage OAuth2 clients for their own account

	ScopeServerCreate         = registerNonServerScope("server.create")
	ScopeServerView           = registerServerScope("server.view")
	ScopeServerAdmin          = registerServerScope("server.admin")
	ScopeServerDelete         = registerServerScope("server.delete")
	ScopeServerEditDefinition = registerServerScope("server.definition.edit")
	ScopeServerEditData       = registerServerScope("server.data.edit")
	ScopeServerEditDataAdmin  = registerServerScope("server.data.edit.admin")
	ScopeServerEditFlags      = registerServerScope("server.flags.edit")
	ScopeServerEditName       = registerServerScope("server.name.edit")
	ScopeServerViewDefinition = registerServerScope("server.definition.view")
	ScopeServerViewData       = registerServerScope("server.data.view")
	ScopeServerViewFlags      = registerServerScope("server.flags.view")

	ScopeServerClientView    = registerServerScope("server.clients.view")
	ScopeServerClientEdit    = registerServerScope("server.clients.edit")
	ScopeServerClientCreate  = registerServerScope("server.clients.create")
	ScopeServerClientDelete  = registerServerScope("server.clients.delete")
	ScopeServerUserView      = registerServerScope("server.users.view")
	ScopeServerUserCreate    = registerServerScope("server.users.create")
	ScopeServerUserEdit      = registerServerScope("server.users.edit")
	ScopeServerUserDelete    = registerServerScope("server.users.delete")
	ScopeServerTaskView      = registerServerScope("server.tasks.view")
	ScopeServerTaskRun       = registerServerScope("server.tasks.run")
	ScopeServerTaskCreate    = registerServerScope("server.tasks.create")
	ScopeServerTaskDelete    = registerServerScope("server.tasks.delete")
	ScopeServerTaskEdit      = registerServerScope("server.tasks.edit")
	ScopeServerReload        = registerServerScope("server.reload")
	ScopeServerStart         = registerServerScope("server.start")
	ScopeServerStop          = registerServerScope("server.stop")
	ScopeServerKill          = registerServerScope("server.kill")
	ScopeServerInstall       = registerServerScope("server.install")
	ScopeServerFileView      = registerServerScope("server.files.view")
	ScopeServerFileEdit      = registerServerScope("server.files.edit")
	ScopeServerSftp          = registerServerScope("server.sftp")
	ScopeServerConsole       = registerServerScope("server.console")
	ScopeServerSendCommand   = registerServerScope("server.console.send")
	ScopeServerStats         = registerServerScope("server.stats")
	ScopeServerStatus        = registerServerScope("server.status")
	ScopeServerBackupView    = registerServerScope("server.backup.view")
	ScopeServerBackupCreate  = registerServerScope("server.backup.create")
	ScopeServerBackupRestore = registerServerScope("server.backup.restore")
	ScopeServerBackupDelete  = registerServerScope("server.backup.delete")

	ScopeSettingsEdit = registerNonServerScope("settings.edit")

	ScopeTemplatesView       = registerNonServerScope("templates.view")
	ScopeTemplatesLocalEdit  = registerNonServerScope("templates.local.edit")
	ScopeTemplatesRepoCreate = registerNonServerScope("templates.repo.create")
	ScopeTemplatesRepoDelete = registerNonServerScope("templates.repo.delete")

	ScopeUserInfoSearch = registerNonServerScope("users.info.search")
	ScopeUserInfoView   = registerNonServerScope("users.info.view")
	ScopeUserInfoEdit   = registerNonServerScope("users.info.edit")
	ScopeUserPermsView  = registerNonServerScope("users.perms.view")
	ScopeUserPermsEdit  = registerNonServerScope("users.perms.edit")

	ScopePanel = registerNonServerScope("panel")
)

func (s *Scope) String() string {
	return s.Value
}

func (s *Scope) Is(t any) bool {
	switch z := t.(type) {
	case string:
		return s.Value == z
	case *Scope:
		return s.Value == z.Value
	default:
		return false
	}
}

func (s *Scope) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Value)
}

func (s *Scope) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	existing := GetScope(str)
	s.Value = existing.Value
	s.ForServer = existing.ForServer
	return nil
}

var allScopes []*Scope

func registerScope(s *Scope) *Scope {
	allScopes = append(allScopes, s)
	return s
}
func registerNonServerScope(s string) *Scope {
	return registerScope(&Scope{Value: s})
}
func registerServerScope(s string) *Scope {
	return registerScope(&Scope{Value: s, ForServer: true})
}

func GetScope(str string) *Scope {
	for _, v := range allScopes {
		if v.Is(str) {
			return v
		}
	}
	return &Scope{Value: str}
}

func ContainsScope(arr []*Scope, value *Scope) bool {
	desired := []*Scope{value}
	if !value.Is(ScopeAdmin.Value) {
		desired = append(desired, ScopeAdmin)
	}
	if value.ForServer && !value.Is(ScopeServerAdmin.Value) {
		desired = append(desired, ScopeServerAdmin)
	}

	for _, v := range arr {
		for _, z := range desired {
			if v.Is(z) {
				return true
			}
		}
	}

	return false
}

func AddScope(source []*Scope, addition *Scope) []*Scope {
	for _, v := range source {
		if v.Is(addition) {
			return source
		}
	}
	return append(source, addition)
}

func RemoveScope(source []*Scope, removal *Scope) []*Scope {
	replacement := make([]*Scope, 0)
	for _, v := range source {
		if !v.Is(removal) {
			replacement = append(replacement, v)
		}
	}
	return replacement
}

func UpdateScopesWhereGranted(source, desired, changer []*Scope) []*Scope {
	replacement := make([]*Scope, 0)
	for _, v := range source {
		//does our user have permission to this scope
		//if so, we need to set this to match the view model
		if ContainsScope(changer, v) {
			if ContainsScope(desired, v) {
				replacement = append(replacement, v)
			}
		} else {
			//otherwise, our current user can't change this value, so re-copy
			replacement = append(replacement, v)
		}
	}
	for _, v := range desired {
		if !ContainsScope(changer, v) {
			continue
		}
		needsAdding := true
		for _, z := range replacement {
			if v.Is(z) {
				needsAdding = false
				break
			}
		}
		if needsAdding {
			replacement = append(replacement, v)
		}
	}
	return replacement
}
