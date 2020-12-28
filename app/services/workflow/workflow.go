package workflow

import (
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/app/services/steps"
	"gitlab.com/scalablespace/listener/lib/components"
)

type workflow struct {
	createNodeByInstanceIdAndPlanId      *steps.CreateNode
	updateRedisPassword                  *steps.UpdateRedisPassword
	migrateOldToNew                      *steps.MigrateRedis
	upgradeDNSRecordByNodeIdAndNewNodeId *steps.UpdateDNSRecord
	setupRedis                           *steps.SetupRedisNode
	setupRedisAsSlave                    *steps.SetupRedisNodeAsSlave
	createDNSRecord                      *steps.CreateDNSRecord
	setupMonitoringByInstanceIdAndNodeId *steps.SetupMonitoring
	finishSetupForInstanceId             *steps.FinishSetup
	disableMonitoringByNodeId            *steps.DisableMonitoring
	destroyDNSRecordByNodeId             *steps.DestroyDNSRecord
	destroyNodeById                      *steps.DestroyNode
	allowAccessFromMultipleIP            *steps.AllowAccessFromMultipleIP
	updateAccessFromWhitelist            *steps.UpdateAccessFromInstanceWhitelist
	finishPoweroff                       *steps.FinishPoweroff
	downloadLatestRedisBackup            *steps.DownloadLatestRedisBackup
	promoteRedisCompatToMaster           *steps.PromotePgToMaster
	createTwoNodes                       *steps.CreateTwoNodes
	setupKeyDBNode                       *steps.SetupKeyDBNode
	setupKeyDBSlaveNode                  *steps.SetupKeydbNodeAsSlave
	setupKeyDBProNode                    *steps.SetupKeyDBProNode
	setupPgNode                          *steps.SetupPgNode
	preparePgWalG                        *steps.PreparePgWalG
	setupPgOnlyPrimaryFromLatestBackup   *steps.SetupPgOnlyPrimaryFromLatestBackup
	setupPgPrimaryFromLatestBackup       *steps.SetupPgPrimaryFromLatestBackup
	setupMySQLNode                       *steps.SetupMySQLNode
	updatePgPassword                     *steps.UpdatePgPassword
	updateMySQLPassword                  *steps.UpdateMySQLPassword
	updateRedisCompatEviction            *steps.UpdateRedisCompatEviction
	setupPgp                             *steps.SetupPGP
	uploadPublicPGP                      *steps.UploadPublicPGP
	prepareMySQLNode                     *steps.PrepareMySQLNode
	setupMySQLNodeFromBackup             *steps.SetupMySQLNodeFromBackup
	setupEncryption                      *steps.SetupEncryption
	pgConfig                             *steps.PgTune
	preparePgWalGWithoutWalPush          *steps.PrepagePgWalGWithoutWalPush
	setupPgReplicaFromLatestBackup       *steps.SetupPgReplicaFromLatestBackup
	promotePgToMaster                    *steps.PromotePgToMaster
	setupPgPrimary                       *steps.SetupPgPrimary
	setupPgNextPrimaryFromLatestBackup   *steps.SetupPgNextPrimaryFromLatestBackup
	promotePgToSingleMaster              *steps.PromotePgToSingleMaster
}

func (w *workflow) Flow(task *models.Task) models.Steps {
	return w.flowForRedis(task)
}

func New(
	createNode *steps.CreateNode,
	updateRedisPassword *steps.UpdateRedisPassword,
	migrateOldToNew *steps.MigrateRedis,
	upgradeDNSRecord *steps.UpdateDNSRecord,
	setupRedis *steps.SetupRedisNode,
	setupRedisAsSlave *steps.SetupRedisNodeAsSlave,
	createDNSRecord *steps.CreateDNSRecord,
	setupMonitoring *steps.SetupMonitoring,
	finishSetup *steps.FinishSetup,
	disableMonitoring *steps.DisableMonitoring,
	destroyDNSRecord *steps.DestroyDNSRecord,
	destroyNode *steps.DestroyNode,
	allowAccessFromMultipleIP *steps.AllowAccessFromMultipleIP,
	allowAccessFromSingleIP *steps.UpdateAccessFromInstanceWhitelist,
	finishPoweroff *steps.FinishPoweroff,
	downloadLatestRedisBackup *steps.DownloadLatestRedisBackup,
	promoteToMaster *steps.PromotePgToMaster,
	createTwoNodes *steps.CreateTwoNodes,
	setupKeyDBNode *steps.SetupKeyDBNode,
	setupKeyDBSlaveNode *steps.SetupKeydbNodeAsSlave,
	setupKeyDBProNode *steps.SetupKeyDBProNode,
	setupPgNode *steps.SetupPgNode,
	prepagePgWalG *steps.PreparePgWalG,
	setupMySQLNode *steps.SetupMySQLNode,
	updatePgPassword *steps.UpdatePgPassword,
	updateMySQLPassword *steps.UpdateMySQLPassword,
	updateRedisCompatEviction *steps.UpdateRedisCompatEviction,
	setupPgp *steps.SetupPGP,
	uploadPublicPGP *steps.UploadPublicPGP,
	prepareMySQLNode *steps.PrepareMySQLNode,
	setupMySQLNodeFromBackup *steps.SetupMySQLNodeFromBackup,
	setupEncryption *steps.SetupEncryption,
	pgConfig *steps.PgTune,
	setupPgOnlyPrimaryFromLatestBackup *steps.SetupPgOnlyPrimaryFromLatestBackup,
	setupPgPrimaryFromLatestBackup *steps.SetupPgPrimaryFromLatestBackup,
	prepagePgWalGWithoutWalPush *steps.PrepagePgWalGWithoutWalPush,
	setupPgReplicaFromLatestBackup *steps.SetupPgReplicaFromLatestBackup,
	promotePgToMaster *steps.PromotePgToMaster,
	setupPgPrimary *steps.SetupPgPrimary,
	setupPgNextPrimaryFromLatestBackup *steps.SetupPgNextPrimaryFromLatestBackup,
	promotePgToSingleMaster *steps.PromotePgToSingleMaster,
) components.Workflow {
	return &workflow{
		createNode,
		updateRedisPassword,
		migrateOldToNew,
		upgradeDNSRecord,
		setupRedis,
		setupRedisAsSlave,
		createDNSRecord,
		setupMonitoring,
		finishSetup,
		disableMonitoring,
		destroyDNSRecord,
		destroyNode,
		allowAccessFromMultipleIP,
		allowAccessFromSingleIP,
		finishPoweroff,
		downloadLatestRedisBackup,
		promoteToMaster,
		createTwoNodes,
		setupKeyDBNode,
		setupKeyDBSlaveNode,
		setupKeyDBProNode,
		setupPgNode,
		prepagePgWalG,
		setupPgOnlyPrimaryFromLatestBackup,
		setupPgPrimaryFromLatestBackup,
		setupMySQLNode,
		updatePgPassword,
		updateMySQLPassword,
		updateRedisCompatEviction,
		setupPgp,
		uploadPublicPGP,
		prepareMySQLNode,
		setupMySQLNodeFromBackup,
		setupEncryption,
		pgConfig,
		prepagePgWalGWithoutWalPush,
		setupPgReplicaFromLatestBackup,
		promotePgToMaster,
		setupPgPrimary,
		setupPgNextPrimaryFromLatestBackup,
		promotePgToSingleMaster,
	}
}
