<template>
  <AppLayout>
    <div class="mx-auto max-w-3xl space-y-6">
      <div class="card">
        <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ t("oauthProvider.adminTitle") }}
          </h2>
          <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
            {{ t("oauthProvider.adminDescription") }}
          </p>
        </div>

        <div v-if="loading" class="flex items-center gap-2 p-6 text-gray-500">
          <div class="h-4 w-4 animate-spin rounded-full border-b-2 border-primary-600"></div>
          {{ t("common.loading") }}
        </div>

        <form v-else class="space-y-5 p-6" @submit.prevent="save">
          <div class="flex items-center justify-between">
            <div>
              <label class="font-medium text-gray-900 dark:text-white">{{ t("oauthProvider.enabled") }}</label>
              <p class="text-sm text-gray-500 dark:text-gray-400">{{ t("oauthProvider.enabledHint") }}</p>
            </div>
            <Toggle v-model="form.enabled" />
          </div>

          <div>
            <label class="input-label">{{ t("oauthProvider.appName") }}</label>
            <input v-model.trim="form.app_name" class="input" />
            <p class="input-hint">{{ t("oauthProvider.appNameHint") }}</p>
          </div>

          <div>
            <label class="input-label">{{ t("oauthProvider.appUrl") }}</label>
            <input v-model.trim="form.app_url" class="input" placeholder="https://video.forapi.cloud" />
            <p class="input-hint">{{ t("oauthProvider.appUrlHint") }}</p>
          </div>

          <div>
            <label class="input-label">{{ t("oauthProvider.clientId") }}</label>
            <input v-model.trim="form.client_id" class="input font-mono" />
          </div>

          <div>
            <label class="input-label">{{ t("oauthProvider.clientSecret") }}</label>
            <div class="flex gap-2">
              <input
                v-model.trim="form.client_secret"
                class="input font-mono"
                autocomplete="off"
                :placeholder="hasSecret ? t('oauthProvider.clientSecretKeep') : ''"
              />
              <button type="button" class="btn btn-secondary flex-shrink-0" @click="generateSecret">
                {{ t("oauthProvider.generate") }}
              </button>
            </div>
            <p class="input-hint">{{ t("oauthProvider.clientSecretHint") }}</p>
          </div>

          <div>
            <label class="input-label">{{ t("oauthProvider.redirectUri") }}</label>
            <input v-model.trim="form.redirect_uri" class="input font-mono" />
            <p class="input-hint">{{ t("oauthProvider.redirectUriHint") }}</p>
          </div>

          <div>
            <label class="input-label">{{ t("oauthProvider.keyQuota") }}</label>
            <input v-model.number="form.key_quota" type="number" min="0" step="0.01" class="input w-40" />
            <p class="input-hint">{{ t("oauthProvider.keyQuotaHint") }}</p>
            <p class="input-hint" data-test="studio-key-groups-hint">{{ t("oauthProvider.keyGroupsAuto") }}</p>
          </div>

          <div class="flex justify-end">
            <button type="submit" class="btn btn-primary" :disabled="saving">
              {{ t("oauthProvider.save") }}
            </button>
          </div>
        </form>
      </div>

      <div class="card p-6">
        <h3 class="font-medium text-gray-900 dark:text-white">{{ t("oauthProvider.endpointsTitle") }}</h3>
        <dl class="mt-3 space-y-2 text-sm">
          <div v-for="item in endpoints" :key="item.label" class="flex flex-wrap gap-2">
            <dt class="w-28 text-gray-500 dark:text-gray-400">{{ item.label }}</dt>
            <dd class="font-mono text-gray-900 dark:text-gray-100">{{ item.url }}</dd>
          </div>
        </dl>
        <p class="mt-3 text-sm text-gray-500 dark:text-gray-400">{{ t("oauthProvider.fieldMapping") }}</p>
      </div>
    </div>
    <TotpStepUpDialog :controller="stepUp" />
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { useI18n } from "vue-i18n";
import oauthAPI, { type OAuthProviderConfig } from "@/api/oauthProvider";
import { useAppStore } from "@/stores";
import AppLayout from "@/components/layout/AppLayout.vue";
import Toggle from "@/components/common/Toggle.vue";
import TotpStepUpDialog from "@/components/auth/TotpStepUpDialog.vue";
import {
  isStepUpBlocked,
  isStepUpCancelled,
  stepUpBlockReason,
  useStepUp,
} from "@/composables/useStepUp";

const { t } = useI18n();
const appStore = useAppStore();
const stepUp = useStepUp();

const loading = ref(true);
const saving = ref(false);
const hasSecret = ref(false);
const form = reactive<OAuthProviderConfig>({
  enabled: false,
  client_id: "",
  client_secret: "",
  redirect_uri: "",
  app_name: "",
  app_url: "",
  // 已弃用：专用密钥改为自动模式；保留字段只为原样回传旧配置
  key_group_id: null,
  key_quota: 0,
});

const endpoints = computed(() => {
  const origin = window.location.origin;
  return [
    { label: t("oauthProvider.authorizationUrl"), url: `${origin}/oauth/authorize` },
    { label: t("oauthProvider.tokenUrl"), url: `${origin}/oauth/token` },
    { label: t("oauthProvider.userInfoUrl"), url: `${origin}/oauth/userinfo` },
  ];
});

function apply(config: OAuthProviderConfig) {
  Object.assign(form, config, { client_secret: "" });
  hasSecret.value = Boolean(config.has_client_secret);
}

function generateSecret() {
  const bytes = new Uint8Array(32);
  crypto.getRandomValues(bytes);
  form.client_secret = Array.from(bytes, (b) => b.toString(16).padStart(2, "0")).join("");
}

function errorMessage(error: unknown): string {
  if (typeof error === "object" && error !== null && "message" in error) {
    return String((error as { message?: unknown }).message || t("common.unknownError"));
  }
  return t("common.unknownError");
}

async function save() {
  saving.value = true;
  try {
    const updated = await stepUp.run(() => oauthAPI.saveAdminConfig({ ...form }));
    apply(updated);
    appStore.showSuccess(t("oauthProvider.saved"));
  } catch (error) {
    if (isStepUpCancelled(error)) return;
    if (isStepUpBlocked(error)) {
      appStore.showError(
        stepUpBlockReason(error) === "STEP_UP_ADMIN_API_KEY_FORBIDDEN"
          ? t("stepUp.adminApiKeyForbidden")
          : t("stepUp.notEnabled"),
      );
      return;
    }
    appStore.showError(errorMessage(error));
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  try {
    apply(await oauthAPI.getAdminConfig());
  } catch (error) {
    appStore.showError(errorMessage(error));
  } finally {
    loading.value = false;
  }
});
</script>
