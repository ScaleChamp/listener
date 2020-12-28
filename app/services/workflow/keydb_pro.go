package workflow

import "gitlab.com/scalablespace/listener/app/models"

func (w *workflow) flowForKeyDBPro(task *models.Task) models.Steps {
	switch task.Action {
	case "create-keydb-pro-1":
		return models.Steps{
			models.NewStep("setup-pgp", w.setupPgp, []string{"instance_id"}),
			models.NewStep("create-vm", w.createNodeByInstanceIdAndPlanId, []string{"instance_id", "plan_id"}, []string{"node_id"}),
			models.NewStep("upload-pgp", w.uploadPublicPGP, []string{"instance_id", "node_id"}),
			models.NewStep("setup-encryption", w.setupEncryption, []string{"node_id"}),
			models.NewStep("update-whitelist", w.updateAccessFromWhitelist, []string{"instance_id"}),
			models.NewStep("setup-vm", w.setupKeyDBProNode, []string{"instance_id", "node_id"}),
			models.NewStep("create-dns", w.createDNSRecord, []string{"node_id", "instance_id"}),
			models.NewStep("setup-monitoring", w.setupMonitoringByInstanceIdAndNodeId, []string{"instance_id", "node_id"}),
			models.NewStep("finish", w.finishSetupForInstanceId, []string{"instance_id"}),
		}
	case "poweron-keydb-pro-1":
		return models.Steps{
			models.NewStep("create-vm", w.createNodeByInstanceIdAndPlanId, []string{"instance_id", "plan_id"}, []string{"node_id"}),
			models.NewStep("upload-pgp", w.uploadPublicPGP, []string{"instance_id", "node_id"}),
			models.NewStep("setup-encryption", w.setupEncryption, []string{"node_id"}),
			models.NewStep("update-whitelist", w.updateAccessFromWhitelist, []string{"instance_id"}),
			models.NewStep("download-backup", w.downloadLatestRedisBackup, []string{"instance_id", "node_id"}),
			models.NewStep("setup-vm", w.setupKeyDBProNode, []string{"instance_id", "node_id"}),
			models.NewStep("create-dns", w.createDNSRecord, []string{"node_id", "instance_id"}),
			models.NewStep("setup-monitoring", w.setupMonitoringByInstanceIdAndNodeId, []string{"instance_id", "node_id"}),
			models.NewStep("finish", w.finishSetupForInstanceId, []string{"instance_id"}),
		}
	case "poweroff-keydb-pro-1":
		return models.Steps{
			models.NewStep("destroy-monitor-1", w.disableMonitoringByNodeId, []string{"node1_id"}),
			models.NewStep("destroy-dns-1", w.destroyDNSRecordByNodeId, []string{"node1_id"}),
			models.NewStep("destroy-node-1", w.destroyNodeById, []string{"node1_id"}),
			models.NewStep("finish-poweroff-1", w.finishPoweroff, []string{"instance_id"}),
		}
	case "destroy-keydb-pro-1":
		return models.Steps{
			models.NewStep("destroy-monitor-1", w.disableMonitoringByNodeId, []string{"node1_id"}),
			models.NewStep("destroy-dns-1", w.destroyDNSRecordByNodeId, []string{"node1_id"}),
			models.NewStep("destroy-node-1", w.destroyNodeById, []string{"node1_id"}),
		}
	case "update-keydb-pro-1-to-1":
		var commands models.Steps
		if task.Metadata["new_password"] != nil {
			commands = append(commands,
				models.NewStep("update-password", w.updateRedisPassword, []string{"instance_id", "node1_id", "old_password"}),
			)
		}
		if task.Metadata["eviction_policy"] != nil {
			commands = append(commands,
				models.NewStep("update-redis-1-eviction-policy", w.updateRedisCompatEviction, []string{"instance_id", "node1_id"}),
			)
		}
		if task.Metadata["plan_id"] != nil {
			commands = append(commands,
				models.NewStep("destroy-monitor-1", w.disableMonitoringByNodeId, []string{"node1_id"}),
				models.NewStep("create-new-vm", w.createNodeByInstanceIdAndPlanId, []string{"instance_id", "plan_id"}, []string{"new_node_id"}),
				models.NewStep("upload-pgp", w.uploadPublicPGP, []string{"instance_id", "new_node_id"}),
				models.NewStep("setup-encryption", w.setupEncryption, []string{"new_node_id"}),
				models.NewStep("setup-ufw", w.allowAccessFromMultipleIP, []string{"instance_id", "node1_id", "new_node_id"}), // updating whitelist for new instance
				models.NewStep("setup-new-vm", w.setupKeyDBProNode, []string{"instance_id", "new_node_id"}),
				models.NewStep("migrate-old-to-new", w.migrateOldToNew, []string{"instance_id", "node1_id", "new_node_id", "password"}),
				models.NewStep("update-dns-record", w.upgradeDNSRecordByNodeIdAndNewNodeId, []string{"instance_id", "node1_id", "new_node_id"}),
				models.NewStep("destroy-old", w.destroyNodeById, []string{"node1_id"}),
				models.NewStep("setup-monitoring-for-1", w.setupMonitoringByInstanceIdAndNodeId, []string{"instance_id", "new_node_id"}),
			)
		}
		if task.Metadata["whitelist"] != nil {
			commands = append(commands,
				models.NewStep("update-whitelist", w.updateAccessFromWhitelist, []string{"instance_id"}), // get whitelist + nodes whitelist + ufw reload
			)
		}
		return append(commands, models.NewStep("finish", w.finishSetupForInstanceId, []string{"instance_id"}))
	case "recover-keydb-pro-master-false":
		return models.Steps{
			models.NewStep("destroy-monitor-1", w.disableMonitoringByNodeId, []string{"node0_id"}),
			models.NewStep("destroy-dns-1", w.destroyDNSRecordByNodeId, []string{"node0_id"}),
			models.NewStep("destroy-node-1", w.destroyNodeById, []string{"node0_id"}),
			models.NewStep("create-vm", w.createNodeByInstanceIdAndPlanId, []string{"instance_id", "plan_id"}, []string{"new_node_id"}),
			models.NewStep("upload-pgp", w.uploadPublicPGP, []string{"instance_id", "new_node_id"}),
			models.NewStep("setup-encryption", w.setupEncryption, []string{"new_node_id"}),
			models.NewStep("update-whitelist", w.updateAccessFromWhitelist, []string{"instance_id"}),
			models.NewStep("download-backup", w.downloadLatestRedisBackup, []string{"instance_id", "new_node_id"}),
			models.NewStep("setup-vm", w.setupKeyDBProNode, []string{"instance_id", "new_node_id"}),
			models.NewStep("create-dns", w.createDNSRecord, []string{"new_node_id", "instance_id"}),
			models.NewStep("setup-monitoring", w.setupMonitoringByInstanceIdAndNodeId, []string{"instance_id", "new_node_id"}),
			models.NewStep("finish", w.finishSetupForInstanceId, []string{"instance_id"}),
		}
	default:
		return w.flowForPostgres(task)
	}
}
