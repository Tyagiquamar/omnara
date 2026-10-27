package maintenance

import (
	"testing"

	"github.com/omnara-ai/omnara/internal/storage"
	"github.com/omnara-ai/omnara/internal/storage/identitystore"
	"github.com/stretchr/testify/require"
)

func TestRunCoreMaintenanceRecoversEachTaskPanic(t *testing.T) {
	result := RunCore(t.Context(), &storage.Store{}, true)
	for _, err := range []error{
		result.ReapRuntimeLocksErr, result.ExpireDaemonRuntimesErr, result.ExpireProcessToolsErr,
		result.WebhookCleanupErr, result.AuthCleanupErr,
		result.IntegrationStatesCleanupErr, result.CompletedInboxCleanupErr, result.DeletedInboxCleanupErr,
	} {
		require.ErrorContains(t, err, "panicked:")
		require.ErrorContains(t, err, "goroutine")
	}
}

func TestRunCoreMaintenanceSkipsCleanupWhenNotDue(t *testing.T) {
	result := RunCore(t.Context(), &storage.Store{}, false)
	for _, err := range []error{
		result.ReapRuntimeLocksErr, result.ExpireDaemonRuntimesErr, result.ExpireProcessToolsErr, result.WebhookCleanupErr,
	} {
		require.ErrorContains(t, err, "panicked:")
	}
	for _, err := range []error{
		result.AuthCleanupErr,
		result.IntegrationStatesCleanupErr, result.CompletedInboxCleanupErr, result.DeletedInboxCleanupErr,
	} {
		require.NoError(t, err)
	}
}

func TestCoreResultCleanupBacklogged(t *testing.T) {
	require.False(t, CoreResult{DeletedWebhooks: int64(eventWebhookCleanupBatchSize), DeletedIntegrationStates: 1}.
		CleanupBacklogged())
	for _, result := range []CoreResult{
		{AuthCleanup: identitystore.AuthStateCleanupResult{DeletedBrowserSessions: 500}},
		{DeletedIntegrationStates: integrationInboxCleanupBatch},
		{CompletedInboxBudgetExhausted: true},
		{DeletedInboxBudgetExhausted: true},
	} {
		require.True(t, result.CleanupBacklogged())
	}
}
