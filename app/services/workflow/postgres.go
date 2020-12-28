package workflow

import "gitlab.com/scalablespace/listener/app/models"

func (w *workflow) flowForPostgres(task *models.Task) models.Steps {
	switch task.Action {
	case "create-pg-1":
		return models.Steps{
			models.NewStep("setup-pgp", w.setupPgp, []string{"instance_id"}),
			models.NewStep("create-vm", w.createNodeByInstanceIdAndPlanId, []string{"instance_id", "plan_id"}, []string{"node_id"}),
			models.NewStep("upload-pgp", w.uploadPublicPGP, []string{"instance_id", "node_id"}),
			models.NewStep("setup-encryption", w.setupEncryption, []string{"node_id"}),
			models.NewStep("update-whitelist", w.updateAccessFromWhitelist, []string{"instance_id"}),
			models.NewStep("update-conf", w.pgConfig, []string{"instance_id", "node_id"}),
			models.NewStep("prepare-vm", w.preparePgWalG, []string{"instance_id", "node_id"}),
			models.NewStep("setup-vm", w.setupPgNode, []string{"instance_id", "node_id"}),
			models.NewStep("create-dns", w.createDNSRecord, []string{"node_id", "instance_id"}),
			models.NewStep("setup-monitoring", w.setupMonitoringByInstanceIdAndNodeId, []string{"instance_id", "node_id"}),
			models.NewStep("finish", w.finishSetupForInstanceId, []string{"instance_id"}),
		}
	case "create-pg-2":
		return models.Steps{
			models.NewStep("setup-pgp", w.setupPgp, []string{"instance_id"}),
			models.NewStep("create-two-vms", w.createTwoNodes, []string{"instance_id", "plan_id"}, []string{"node1_id", "node2_id"}),
			models.NewStep("upload-pgp-1", w.uploadPublicPGP, []string{"instance_id", "node1_id"}),
			models.NewStep("upload-pgp-2", w.uploadPublicPGP, []string{"instance_id", "node2_id"}),
			models.NewStep("setup-encryption-1", w.setupEncryption, []string{"node1_id"}),
			models.NewStep("setup-encryption-2", w.setupEncryption, []string{"node2_id"}),
			models.NewStep("setup-ufw-for-first-vm", w.allowAccessFromMultipleIP, []string{"instance_id", "node1_id", "node2_id"}),
			models.NewStep("setup-ufw-for-second-vm", w.allowAccessFromMultipleIP, []string{"instance_id", "node2_id"}),
			// setup first vm
			models.NewStep("update-conf-1", w.pgConfig, []string{"instance_id", "node1_id"}),
			models.NewStep("prepare-vm-1", w.preparePgWalG, []string{"instance_id", "node1_id"}),
			models.NewStep("setup-vm-1", w.setupPgPrimary, []string{"instance_id", "node1_id"}), // withbackup
			// setup second vm
			models.NewStep("update-conf-2", w.pgConfig, []string{"instance_id", "node2_id"}),
			models.NewStep("prepare-vm-2", w.preparePgWalG, []string{"instance_id", "node2_id"}),
			models.NewStep("setup-vm-2", w.setupPgReplicaFromLatestBackup, []string{"instance_id", "node1_id", "node2_id"}),

			models.NewStep("create-dns-for-first-vm", w.createDNSRecord, []string{"node1_id", "instance_id"}),
			models.NewStep("create-dns-for-second-vm", w.createDNSRecord, []string{"node2_id", "instance_id"}),
			models.NewStep("setup-monitoring-for-1", w.setupMonitoringByInstanceIdAndNodeId, []string{"instance_id", "node1_id"}),
			models.NewStep("setup-monitoring-for-2", w.setupMonitoringByInstanceIdAndNodeId, []string{"instance_id", "node2_id"}),
			models.NewStep("finish", w.finishSetupForInstanceId, []string{"instance_id"}),
		}
	case "poweron-pg-2":
		return models.Steps{
			models.NewStep("create-two-vms", w.createTwoNodes, []string{"instance_id", "plan_id"}, []string{"node1_id", "node2_id"}),
			models.NewStep("upload-pgp-1", w.uploadPublicPGP, []string{"instance_id", "node1_id"}),
			models.NewStep("upload-pgp-2", w.uploadPublicPGP, []string{"instance_id", "node2_id"}),
			models.NewStep("setup-encryption-1", w.setupEncryption, []string{"node1_id"}),
			models.NewStep("setup-encryption-2", w.setupEncryption, []string{"node2_id"}),
			models.NewStep("setup-ufw-for-first-vm", w.allowAccessFromMultipleIP, []string{"instance_id", "node1_id", "node2_id"}),
			models.NewStep("setup-ufw-for-second-vm", w.allowAccessFromMultipleIP, []string{"instance_id", "node2_id"}),
			models.NewStep("update-conf-1", w.pgConfig, []string{"instance_id", "node1_id"}),
			models.NewStep("prepare-vm-1", w.preparePgWalG, []string{"instance_id", "node1_id"}),
			models.NewStep("update-conf-2", w.pgConfig, []string{"instance_id", "node2_id"}),
			models.NewStep("prepare-vm-2", w.preparePgWalG, []string{"instance_id", "node2_id"}),
			models.NewStep("setup-vm-1", w.setupPgPrimaryFromLatestBackup, []string{"instance_id", "node1_id"}),
			models.NewStep("setup-vm-2", w.setupPgReplicaFromLatestBackup, []string{"instance_id", "node1_id", "node2_id"}),
			models.NewStep("create-dns-for-first-vm", w.createDNSRecord, []string{"node1_id", "instance_id"}),
			models.NewStep("create-dns-for-second-vm", w.createDNSRecord, []string{"node2_id", "instance_id"}),
			models.NewStep("setup-monitoring-for-1", w.setupMonitoringByInstanceIdAndNodeId, []string{"instance_id", "node1_id"}),
			models.NewStep("setup-monitoring-for-2", w.setupMonitoringByInstanceIdAndNodeId, []string{"instance_id", "node2_id"}),
			models.NewStep("finish", w.finishSetupForInstanceId, []string{"instance_id"}),
		}
	case "poweroff-pg-2":
		return models.Steps{
			models.NewStep("destroy-monitor-1", w.disableMonitoringByNodeId, []string{"node1_id"}),
			models.NewStep("destroy-monitor-2", w.disableMonitoringByNodeId, []string{"node2_id"}),
			models.NewStep("destroy-dns-1", w.destroyDNSRecordByNodeId, []string{"node1_id"}),
			models.NewStep("destroy-dns-2", w.destroyDNSRecordByNodeId, []string{"node2_id"}),
			models.NewStep("destroy-node-1", w.destroyNodeById, []string{"node1_id"}),
			models.NewStep("destroy-node-2", w.destroyNodeById, []string{"node2_id"}),
			models.NewStep("finish-poweroff-1", w.finishPoweroff, []string{"instance_id"}),
		}
	case "poweron-pg-1":
		return models.Steps{
			models.NewStep("create-vm", w.createNodeByInstanceIdAndPlanId, []string{"instance_id", "plan_id"}, []string{"node_id"}),
			models.NewStep("upload-pgp", w.uploadPublicPGP, []string{"instance_id", "node_id"}),
			models.NewStep("setup-encryption", w.setupEncryption, []string{"node_id"}),
			models.NewStep("update-whitelist", w.updateAccessFromWhitelist, []string{"instance_id"}),
			models.NewStep("update-conf", w.pgConfig, []string{"instance_id", "node_id"}),
			models.NewStep("prepare-vm", w.preparePgWalG, []string{"instance_id", "node_id"}),
			models.NewStep("setup-from-latest-backup", w.setupPgOnlyPrimaryFromLatestBackup, []string{"instance_id", "node_id"}),
			models.NewStep("create-dns", w.createDNSRecord, []string{"node_id", "instance_id"}),
			models.NewStep("setup-monitoring", w.setupMonitoringByInstanceIdAndNodeId, []string{"instance_id", "node_id"}),
			models.NewStep("finish", w.finishSetupForInstanceId, []string{"instance_id"}),
		}
	case "poweroff-pg-1":
		return models.Steps{
			models.NewStep("destroy-monitor-1", w.disableMonitoringByNodeId, []string{"node1_id"}),
			models.NewStep("destroy-dns-1", w.destroyDNSRecordByNodeId, []string{"node1_id"}),
			models.NewStep("destroy-node-1", w.destroyNodeById, []string{"node1_id"}),
			models.NewStep("finish-poweroff-1", w.finishPoweroff, []string{"instance_id"}),
		}
	case "update-pg-1-to-2":
		var commands models.Steps
		if task.Metadata["new_password"] != nil {
			commands = append(commands,
				models.NewStep("update-password", w.updatePgPassword, []string{"instance_id", "node1_id"}),
			)
		}
		if task.Metadata["plan_id"] != nil {
			commands = append(commands,
				models.NewStep("destroy-monitor-1", w.disableMonitoringByNodeId, []string{"node1_id"}),
				models.NewStep("create-two-vms", w.createTwoNodes, []string{"instance_id", "plan_id"}, []string{"new_master_id", "new_slave_id"}),
				models.NewStep("upload-pgp-1", w.uploadPublicPGP, []string{"instance_id", "new_master_id"}),
				models.NewStep("upload-pgp-2", w.uploadPublicPGP, []string{"instance_id", "new_slave_id"}),
				models.NewStep("setup-encryption-1", w.setupEncryption, []string{"new_master_id"}),
				models.NewStep("setup-encryption-2", w.setupEncryption, []string{"new_slave_id"}),
				models.NewStep("setup-ufw-for-old-master", w.allowAccessFromMultipleIP, []string{"instance_id", "node1_id", "new_master_id"}),
				models.NewStep("setup-ufw-for-new-master", w.allowAccessFromMultipleIP, []string{"instance_id", "new_master_id", "new_slave_id"}),
				models.NewStep("setup-ufw-for-new-slave", w.allowAccessFromMultipleIP, []string{"instance_id", "new_slave_id"}),
				models.NewStep("update-conf-1", w.pgConfig, []string{"instance_id", "new_master_id"}),
				models.NewStep("update-conf-2", w.pgConfig, []string{"instance_id", "new_slave_id"}),
				models.NewStep("prepare-vm-1", w.preparePgWalG, []string{"instance_id", "new_master_id"}),
				models.NewStep("prepare-vm-2", w.preparePgWalG, []string{"instance_id", "new_slave_id"}),
				models.NewStep("setup-vm-1", w.setupPgNextPrimaryFromLatestBackup, []string{"instance_id", "node1_id", "new_master_id"}),
				models.NewStep("upgrade-dns-1", w.upgradeDNSRecordByNodeIdAndNewNodeId, []string{"instance_id", "node1_id", "new_master_id"}),
				models.NewStep("destroy-old", w.destroyNodeById, []string{"node1_id"}),
				models.NewStep("promote-vm-1", w.promotePgToMaster, []string{"new_master_id"}), // failover proxy should be here then destroy node
				models.NewStep("setup-vm-2", w.setupPgReplicaFromLatestBackup, []string{"instance_id", "new_master_id", "new_slave_id"}),
				models.NewStep("create-dns-for-new-slave", w.createDNSRecord, []string{"new_slave_id", "instance_id"}),
				models.NewStep("setup-monitoring-for-1", w.setupMonitoringByInstanceIdAndNodeId, []string{"instance_id", "new_master_id"}),
				models.NewStep("setup-monitoring-for-2", w.setupMonitoringByInstanceIdAndNodeId, []string{"instance_id", "new_slave_id"}),
			)
		}
		if task.Metadata["whitelist"] != nil {
			commands = append(commands,
				models.NewStep("update-whitelist", w.updateAccessFromWhitelist, []string{"instance_id"}), // get whitelist + nodes whitelist + ufw reload
			)
		}
		return append(commands, models.NewStep("finish", w.finishSetupForInstanceId, []string{"instance_id"}))
	case "update-pg-2-to-1": // TODO
		var commands models.Steps
		if task.Metadata["new_password"] != nil {
			commands = append(commands,
				models.NewStep("update-password", w.updatePgPassword, []string{"instance_id", "node1_id"}),
			)
		}
		if task.Metadata["plan_id"] != nil {
			commands = append(commands,
				models.NewStep("destroy-monitor-1", w.disableMonitoringByNodeId, []string{"node1_id"}),
				models.NewStep("destroy-monitor-2", w.disableMonitoringByNodeId, []string{"node2_id"}),
				models.NewStep("destroy-dns-2", w.destroyDNSRecordByNodeId, []string{"node2_id"}),
				models.NewStep("destroy-node-2", w.destroyNodeById, []string{"node2_id"}),
				models.NewStep("create-new-vm", w.createNodeByInstanceIdAndPlanId, []string{"instance_id", "plan_id"}, []string{"new_master_id"}),
				models.NewStep("upload-pgp", w.uploadPublicPGP, []string{"instance_id", "new_master_id"}),
				models.NewStep("setup-encryption", w.setupEncryption, []string{"new_master_id"}),
				models.NewStep("setup-ufw", w.updateAccessFromWhitelist, []string{"instance_id"}),
				models.NewStep("setup-ufw-for-new-master", w.allowAccessFromMultipleIP, []string{"instance_id", "node1_id", "new_master_id"}),
				models.NewStep("update-conf-2", w.pgConfig, []string{"instance_id", "new_master_id"}),
				models.NewStep("prepare-vm-2", w.preparePgWalG, []string{"instance_id", "new_master_id"}),
				models.NewStep("setup-vm-2", w.setupPgNextPrimaryFromLatestBackup, []string{"instance_id", "node1_id", "new_master_id"}),
				models.NewStep("upgrade-dns-2", w.upgradeDNSRecordByNodeIdAndNewNodeId, []string{"instance_id", "node1_id", "new_master_id"}),
				models.NewStep("destroy-node-1", w.destroyNodeById, []string{"node1_id"}),
				models.NewStep("promote-vm-2", w.promotePgToSingleMaster, []string{"new_master_id"}), // failover proxy should be here
				models.NewStep("setup-monitoring-for-1", w.setupMonitoringByInstanceIdAndNodeId, []string{"instance_id", "new_master_id"}),
			)
		}
		if task.Metadata["whitelist"] != nil {
			commands = append(commands,
				models.NewStep("update-whitelist", w.updateAccessFromWhitelist, []string{"instance_id"}), // get whitelist + nodes whitelist + ufw reload
			)
		}
		return append(commands, models.NewStep("finish", w.finishSetupForInstanceId, []string{"instance_id"}))
	case "update-pg-2-to-2":
		var commands models.Steps
		if task.Metadata["new_password"] != nil {
			commands = append(commands,
				models.NewStep("update-password", w.updatePgPassword, []string{"instance_id", "node1_id"}),
			)
		}
		if task.Metadata["plan_id"] != nil {
			commands = append(commands,
				models.NewStep("destroy-monitor-1", w.disableMonitoringByNodeId, []string{"node1_id"}),
				models.NewStep("destroy-monitor-2", w.disableMonitoringByNodeId, []string{"node2_id"}),
				models.NewStep("create-two-vms", w.createTwoNodes, []string{"instance_id", "plan_id"}, []string{"new_master_id", "new_slave_id"}),
				models.NewStep("upload-pgp-1", w.uploadPublicPGP, []string{"instance_id", "new_master_id"}),
				models.NewStep("upload-pgp-2", w.uploadPublicPGP, []string{"instance_id", "new_slave_id"}),
				models.NewStep("setup-encryption-1", w.setupEncryption, []string{"new_master_id"}),
				models.NewStep("setup-encryption-2", w.setupEncryption, []string{"new_slave_id"}),
				models.NewStep("setup-ufw-for-new-master", w.allowAccessFromMultipleIP, []string{"instance_id", "new_master_id", "new_slave_id"}),
				models.NewStep("setup-ufw-for-old-master", w.allowAccessFromMultipleIP, []string{"instance_id", "node1_id", "node2_id", "new_master_id"}),
				models.NewStep("setup-ufw", w.updateAccessFromWhitelist, []string{"instance_id"}),
				models.NewStep("update-conf-1", w.pgConfig, []string{"instance_id", "new_master_id"}),
				models.NewStep("prepare-vm-1", w.preparePgWalG, []string{"instance_id", "new_master_id"}),
				models.NewStep("setup-vm-1", w.setupPgNextPrimaryFromLatestBackup, []string{"instance_id", "node1_id", "new_master_id"}),
				models.NewStep("upgrade-dns-1", w.upgradeDNSRecordByNodeIdAndNewNodeId, []string{"instance_id", "node1_id", "new_master_id"}),
				models.NewStep("destroy-node-1", w.destroyNodeById, []string{"node1_id"}),
				models.NewStep("promote-vm-1", w.promotePgToMaster, []string{"new_master_id"}),
				models.NewStep("update-conf-2", w.pgConfig, []string{"instance_id", "new_slave_id"}),
				models.NewStep("prepare-vm-2", w.preparePgWalG, []string{"instance_id", "new_slave_id"}),
				models.NewStep("setup-vm-2", w.setupPgReplicaFromLatestBackup, []string{"instance_id", "new_master_id", "new_slave_id"}),
				models.NewStep("upgrade-dns-2", w.upgradeDNSRecordByNodeIdAndNewNodeId, []string{"instance_id", "node2_id", "new_slave_id"}),
				models.NewStep("destroy-node-2", w.destroyNodeById, []string{"node2_id"}),
				models.NewStep("setup-monitoring-for-1", w.setupMonitoringByInstanceIdAndNodeId, []string{"instance_id", "new_master_id"}),
				models.NewStep("setup-monitoring-for-2", w.setupMonitoringByInstanceIdAndNodeId, []string{"instance_id", "new_slave_id"}),
			)
		}
		if task.Metadata["whitelist"] != nil {
			commands = append(commands,
				models.NewStep("update-whitelist", w.updateAccessFromWhitelist, []string{"instance_id"}), // get whitelist + nodes whitelist + ufw reload
			)
		}
		return append(commands, models.NewStep("finish", w.finishSetupForInstanceId, []string{"instance_id"}))
	case "update-pg-1-to-1":
		var commands models.Steps
		if task.Metadata["new_password"] != nil {
			commands = append(commands,
				models.NewStep("update-password", w.updatePgPassword, []string{"instance_id", "node1_id"}),
			)
		}
		if task.Metadata["plan_id"] != nil {
			commands = append(commands,
				models.NewStep("destroy-monitor-1", w.disableMonitoringByNodeId, []string{"node1_id"}),
				models.NewStep("destroy-dns-1", w.destroyDNSRecordByNodeId, []string{"node1_id"}),
				models.NewStep("destroy-node-1", w.destroyNodeById, []string{"node1_id"}),
				models.NewStep("create-vm", w.createNodeByInstanceIdAndPlanId, []string{"instance_id", "plan_id"}, []string{"new_node_id"}),
				models.NewStep("upload-pgp", w.uploadPublicPGP, []string{"instance_id", "new_node_id"}),
				models.NewStep("setup-encryption", w.setupEncryption, []string{"new_node_id"}),
				models.NewStep("update-whitelist", w.updateAccessFromWhitelist, []string{"instance_id"}),
				models.NewStep("update-conf", w.pgConfig, []string{"instance_id", "new_node_id"}),
				models.NewStep("prepare-vm", w.preparePgWalG, []string{"instance_id", "new_node_id"}),
				models.NewStep("setup-from-latest-backup", w.setupPgOnlyPrimaryFromLatestBackup, []string{"instance_id", "new_node_id"}),
				models.NewStep("create-dns", w.createDNSRecord, []string{"new_node_id", "instance_id"}),
				models.NewStep("setup-monitoring", w.setupMonitoringByInstanceIdAndNodeId, []string{"instance_id", "new_node_id"}),
				models.NewStep("finish", w.finishSetupForInstanceId, []string{"instance_id"}),
			)
		}
		if task.Metadata["whitelist"] != nil {
			commands = append(commands,
				models.NewStep("update-whitelist", w.updateAccessFromWhitelist, []string{"instance_id"}), // get whitelist + nodes whitelist + ufw reload
			)
		}
		return append(commands, models.NewStep("finish", w.finishSetupForInstanceId, []string{"instance_id"}))
	case "recover-pg-master-false":
		return models.Steps{
			models.NewStep("destroy-monitor-1", w.disableMonitoringByNodeId, []string{"node0_id"}),
			models.NewStep("destroy-dns-1", w.destroyDNSRecordByNodeId, []string{"node0_id"}),
			models.NewStep("destroy-node-1", w.destroyNodeById, []string{"node0_id"}),
			models.NewStep("create-vm", w.createNodeByInstanceIdAndPlanId, []string{"instance_id", "plan_id"}, []string{"new_node_id"}),
			models.NewStep("upload-pgp", w.uploadPublicPGP, []string{"instance_id", "new_node_id"}),
			models.NewStep("setup-encryption", w.setupEncryption, []string{"new_node_id"}),
			models.NewStep("update-conf", w.pgConfig, []string{"instance_id", "new_node_id"}),
			models.NewStep("update-whitelist", w.updateAccessFromWhitelist, []string{"instance_id"}),
			models.NewStep("prepare-vm", w.preparePgWalG, []string{"instance_id", "new_node_id"}),
			models.NewStep("setup-from-latest-backup", w.setupPgOnlyPrimaryFromLatestBackup, []string{"instance_id", "new_node_id"}),
			models.NewStep("create-dns", w.createDNSRecord, []string{"new_node_id", "instance_id"}),
			models.NewStep("setup-monitoring", w.setupMonitoringByInstanceIdAndNodeId, []string{"instance_id", "new_node_id"}),
			models.NewStep("finish", w.finishSetupForInstanceId, []string{"instance_id"}),
		}
	case "recover-pg-master-false-slave-true":
		return models.Steps{
			models.NewStep("destroy-monitor-for-old-master", w.disableMonitoringByNodeId, []string{"node0_id"}),
			models.NewStep("destroy-dns-for-old-slave", w.destroyDNSRecordByNodeId, []string{"node1_id"}),
			models.NewStep("destroy-old-master", w.destroyNodeById, []string{"node0_id"}),
			models.NewStep("promote-slave-to-master", w.promotePgToMaster, []string{"node1_id"}),
			models.NewStep("update-dns-record-from-old-master-to-new", w.upgradeDNSRecordByNodeIdAndNewNodeId, []string{"instance_id", "node0_id", "node1_id"}),
			models.NewStep("create-new-vm", w.createNodeByInstanceIdAndPlanId, []string{"instance_id", "plan_id"}, []string{"node2_id"}),
			models.NewStep("upload-pgp", w.uploadPublicPGP, []string{"instance_id", "node2_id"}),
			models.NewStep("setup-encryption", w.setupEncryption, []string{"node2_id"}),
			models.NewStep("setup-ufw-for-new-master", w.allowAccessFromMultipleIP, []string{"instance_id", "node1_id", "node2_id"}),
			models.NewStep("setup-ufw-for-new-slave", w.allowAccessFromMultipleIP, []string{"instance_id", "node2_id"}),
			models.NewStep("update-conf-2", w.pgConfig, []string{"instance_id", "node2_id"}),
			models.NewStep("prepare-vm-2", w.preparePgWalG, []string{"instance_id", "node2_id"}),
			models.NewStep("setup-vm-2", w.setupPgReplicaFromLatestBackup, []string{"instance_id", "node1_id", "node2_id"}),
			models.NewStep("create-slave-dns-for-new-slave", w.createDNSRecord, []string{"node2_id", "instance_id"}),
			models.NewStep("setup-monitoring", w.setupMonitoringByInstanceIdAndNodeId, []string{"instance_id", "node2_id"}),
			models.NewStep("finish", w.finishSetupForInstanceId, []string{"instance_id"}),
		}
	case "recover-pg-master-true-slave-false":
		return models.Steps{
			models.NewStep("destroy-node-1", w.destroyNodeById, []string{"node1_id"}),
			models.NewStep("destroy-monitor-1", w.disableMonitoringByNodeId, []string{"node1_id"}),
			models.NewStep("create-new-vm", w.createNodeByInstanceIdAndPlanId, []string{"instance_id", "plan_id"}, []string{"node2_id"}),
			models.NewStep("upload-pgp", w.uploadPublicPGP, []string{"instance_id", "node2_id"}),
			models.NewStep("setup-encryption", w.setupEncryption, []string{"node2_id"}),
			models.NewStep("setup-ufw", w.allowAccessFromMultipleIP, []string{"instance_id", "node0_id", "node2_id"}),
			models.NewStep("update-conf-2", w.pgConfig, []string{"instance_id", "node2_id"}),
			models.NewStep("prepare-vm-2", w.preparePgWalG, []string{"instance_id", "node2_id"}),
			models.NewStep("setup-vm-2", w.setupPgReplicaFromLatestBackup, []string{"instance_id", "node0_id", "node2_id"}),
			models.NewStep("update-dns-record", w.upgradeDNSRecordByNodeIdAndNewNodeId, []string{"instance_id", "node1_id", "node2_id"}),
		}
	case "recover-pg-master-false-slave-false":
		return models.Steps{
			models.NewStep("destroy-monitor-1", w.disableMonitoringByNodeId, []string{"node0_id"}),
			models.NewStep("destroy-monitor-2", w.disableMonitoringByNodeId, []string{"node1_id"}),
			models.NewStep("destroy-dns-1", w.destroyDNSRecordByNodeId, []string{"node0_id"}),
			models.NewStep("destroy-dns-2", w.destroyDNSRecordByNodeId, []string{"node1_id"}),
			models.NewStep("destroy-node-1", w.destroyNodeById, []string{"node0_id"}),
			models.NewStep("destroy-node-2", w.destroyNodeById, []string{"node1_id"}),
			models.NewStep("create-two-vms", w.createTwoNodes, []string{"instance_id", "plan_id"}, []string{"new_node1_id", "new_node2_id"}),
			models.NewStep("upload-pgp-1", w.uploadPublicPGP, []string{"instance_id", "new_node1_id"}),
			models.NewStep("upload-pgp-2", w.uploadPublicPGP, []string{"instance_id", "new_node2_id"}),
			models.NewStep("setup-encryption-1", w.setupEncryption, []string{"new_node1_id"}),
			models.NewStep("setup-encryption-2", w.setupEncryption, []string{"new_node2_id"}),
			models.NewStep("update-whitelist", w.updateAccessFromWhitelist, []string{"instance_id"}),
			models.NewStep("setup-ufw", w.allowAccessFromMultipleIP, []string{"instance_id", "new_node1_id", "new_node2_id"}),
			models.NewStep("update-conf-1", w.pgConfig, []string{"instance_id", "new_node1_id"}),
			models.NewStep("prepare-vm-1", w.preparePgWalG, []string{"instance_id", "new_node1_id"}),
			models.NewStep("setup-vm-1", w.setupPgPrimaryFromLatestBackup, []string{"instance_id", "new_node1_id"}),
			models.NewStep("update-conf-2", w.pgConfig, []string{"instance_id", "new_node2_id"}),
			models.NewStep("prepare-vm-2", w.preparePgWalG, []string{"instance_id", "new_node2_id"}),
			models.NewStep("setup-vm-2", w.setupPgReplicaFromLatestBackup, []string{"instance_id", "new_node1_id", "new_node2_id"}),
			models.NewStep("create-dns-for-first-vm", w.createDNSRecord, []string{"new_node1_id", "instance_id"}),
			models.NewStep("create-dns-for-second-vm", w.createDNSRecord, []string{"new_node2_id", "instance_id"}),
			models.NewStep("setup-monitoring-for-1", w.setupMonitoringByInstanceIdAndNodeId, []string{"instance_id", "new_node1_id"}),
			models.NewStep("setup-monitoring-for-2", w.setupMonitoringByInstanceIdAndNodeId, []string{"instance_id", "new_node2_id"}),
			models.NewStep("finish", w.finishSetupForInstanceId, []string{"instance_id"}),
		}
	case "destroy-pg-1":
		return models.Steps{
			models.NewStep("destroy-monitor-1", w.disableMonitoringByNodeId, []string{"node1_id"}),
			models.NewStep("destroy-dns-1", w.destroyDNSRecordByNodeId, []string{"node1_id"}),
			models.NewStep("destroy-node-1", w.destroyNodeById, []string{"node1_id"}),
		}
	case "destroy-pg-2":
		return models.Steps{
			models.NewStep("destroy-monitor-1", w.disableMonitoringByNodeId, []string{"node1_id"}),
			models.NewStep("destroy-monitor-2", w.disableMonitoringByNodeId, []string{"node2_id"}),
			models.NewStep("destroy-dns-1", w.destroyDNSRecordByNodeId, []string{"node1_id"}),
			models.NewStep("destroy-dns-2", w.destroyDNSRecordByNodeId, []string{"node2_id"}),
			models.NewStep("destroy-node-1", w.destroyNodeById, []string{"node1_id"}),
			models.NewStep("destroy-node-2", w.destroyNodeById, []string{"node2_id"}),
		}
	default:
		return w.flowForMySQL(task)
	}
}
