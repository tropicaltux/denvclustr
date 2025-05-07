package schema

// Postprocess applies defaults and fallbacks after validation.
func Postprocess(root *DenvclustrRoot) {
	nodes := collectNodeMap(root)

	for _, devcontainer := range root.Devcontainers {
		// ensure at least one remote‑access mechanism exists
		if devcontainer.RemoteAccess == nil || (devcontainer.RemoteAccess.OpenVsCodeServer == nil && devcontainer.RemoteAccess.Ssh == nil) {
			devcontainer.RemoteAccess = &DevcontainerRemoteAccess{OpenVsCodeServer: &DevcontainerOpenVSCodeServer{}}
		}

		// fallback for SSH public key
		if devcontainer.RemoteAccess.Ssh != nil && devcontainer.RemoteAccess.Ssh.PublicSshKey == "" {
			devcontainer.RemoteAccess.Ssh.PublicSshKey = nodes[string(devcontainer.NodeId)].RemoteAccess.PublicSSHKey
		}
	}
}
