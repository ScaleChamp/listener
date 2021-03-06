package config

import (
	"gitlab.com/scalablespace/listener/app/controllers"
	"gitlab.com/scalablespace/listener/app/services"
	"gitlab.com/scalablespace/listener/app/services/steps"
	"gitlab.com/scalablespace/listener/app/services/workflow"
	"gitlab.com/scalablespace/listener/config/initializers"
	"gitlab.com/scalablespace/listener/db"
	"gitlab.com/scalablespace/listener/db/repositories"
	"gitlab.com/scalablespace/listener/lib/clouds"
	"gitlab.com/scalablespace/listener/lib/clouds/adapters"
	"gitlab.com/scalablespace/listener/lib/executor"
	"gitlab.com/scalablespace/listener/lib/taskflow"
	"go.uber.org/fx"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(newEnvironment),
		fx.Provide(db.NewDB),
		fx.Provide(initializers.NewDeliveries),
		fx.Provide(initializers.NewPusherClient),
		fx.Provide(initializers.NewHTTPClient),
		fx.Provide(workflow.New),
		fx.Provide(taskflow.NewEngine),
		fx.Provide(services.NewConsumer),
		fx.Provide(steps.NewSetupRedisNode),
		fx.Provide(steps.NewUploadPublicPGP),
		fx.Provide(steps.NewSetupPGP),
		fx.Provide(steps.NewSetupPgNode),
		fx.Provide(steps.NewSetupMySQLNode),
		fx.Provide(steps.NewSetupPgPrimaryFromLatestBackup),
		fx.Provide(steps.NewSetupPgPrimary),
		fx.Provide(steps.NewSetupPgOnlyPrimaryFromLatestBackup),
		fx.Provide(steps.NewSetupPgReplicaFromLatestBackup),
		fx.Provide(steps.NewSetupPgNextPrimaryFromLatestBackup),
		fx.Provide(steps.NewPreparePgNode),
		fx.Provide(steps.NewPreparePgWalGWithoutWalPush),
		fx.Provide(steps.NewUpdateAccessFromInstanceWhitelist),
		fx.Provide(steps.NewAllowAccessFromMultipleIP),
		fx.Provide(steps.NewMigrateRedis),
		fx.Provide(steps.NewUpdateDNSRecord),
		fx.Provide(steps.NewUpdatePgPassword),
		fx.Provide(steps.NewUpdateMySQLPassword),
		fx.Provide(steps.NewUpdateRedisPassword),
		fx.Provide(steps.NewUpdateRedisCompatEviction),
		fx.Provide(steps.NewSetupNodeAsSlave),
		fx.Provide(steps.NewSetupKeydbNodeAsSlave),
		fx.Provide(steps.NewSetupKeyDBProNode),
		fx.Provide(steps.NewSetupKeyDBNode),
		fx.Provide(steps.NewSetupEncryption),
		fx.Provide(steps.NewSetupMonitoring),
		fx.Provide(steps.NewPrepareMySQLNode),
		fx.Provide(steps.NewPgTune),
		fx.Provide(steps.NewSetupMySQLNodeFromBackup),
		fx.Provide(steps.NewFinishSetup),
		fx.Provide(steps.NewFinishPoweroff),
		fx.Provide(steps.NewDownloadLatestRedisBackup),
		fx.Provide(steps.NewPromoteToMasterRedisCompat),
		fx.Provide(steps.NewPromotePgToMaster),
		fx.Provide(steps.NewPromotePgToSingleMaster),
		fx.Provide(steps.NewDisableMonitoring),
		fx.Provide(steps.NewDestroyNode),
		fx.Provide(steps.NewDestroyDNSRecord),
		fx.Provide(steps.NewCreateNode),
		fx.Provide(steps.NewCreateTwoNodes),
		fx.Provide(steps.NewCreateDNSRecord),
		fx.Provide(executor.NewExecutor),
		fx.Provide(adapters.NewDO),
		fx.Provide(adapters.NewAzure),
		fx.Provide(adapters.NewGCP),
		fx.Provide(adapters.NewExoscale),
		fx.Provide(adapters.NewScaleway),
		fx.Provide(adapters.NewIBM),
		fx.Provide(adapters.NewAlibaba),
		fx.Provide(adapters.NewSelectel),
		fx.Provide(adapters.NewTencentCloud),
		fx.Provide(adapters.NewLinode),
		fx.Provide(adapters.NewOracle),
		fx.Provide(adapters.NewHetzner),
		fx.Provide(adapters.NewAWS),
		fx.Provide(adapters.NewUpcloud),
		fx.Provide(adapters.NewVultr),
		fx.Provide(clouds.NewCloudAdapter),
		fx.Provide(initializers.NewDOClient),
		fx.Provide(initializers.NewGCPClient),
		fx.Provide(initializers.NewExoscale),
		fx.Provide(initializers.NewIBMClient),
		fx.Provide(initializers.NewLinodeGo),
		fx.Provide(initializers.NewUpcloud),
		fx.Provide(initializers.NewCloudflare),
		fx.Provide(initializers.NewHetznerClient),
		fx.Provide(initializers.NewAzureClient),
		fx.Provide(initializers.NewVultrClient),
		fx.Provide(initializers.NewScalewayClient),
		fx.Provide(controllers.NewHealthController),
		fx.Provide(repositories.NewInstanceRepository),
		fx.Provide(repositories.NewAccessKeyPairRepository),
		fx.Provide(repositories.NewCertificateAuthorityRepository),
		fx.Provide(repositories.NewEncryptionKeyRepository),
		fx.Provide(repositories.NewNodeRepository),
		fx.Provide(repositories.NewTaskRepository),
		fx.Provide(repositories.NewUsageRepository),
		fx.Provide(repositories.NewPlanRepository),
		fx.Provide(repositories.NewPrometheusRepository),
		fx.Provide(routes),
		fx.Invoke(boot),
	)
}
