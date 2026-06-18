<template>
  <div class="row">
    <div class="column">
      <form class="card" @submit="updateSettings">
        <div class="card-title">
          <h2>{{ t("settings.profileSettings") }}</h2>
        </div>

        <div class="card-content">
          <p>
            <input type="checkbox" name="hideDotfiles" v-model="hideDotfiles" />
            {{ t("settings.hideDotfiles") }}
          </p>
          <p>
            <input type="checkbox" name="singleClick" v-model="singleClick" />
            {{ t("settings.singleClick") }}
          </p>
          <p>
            <input
              type="checkbox"
              name="redirectAfterCopyMove"
              v-model="redirectAfterCopyMove"
            />
            {{ t("settings.redirectAfterCopyMove") }}
          </p>
          <p>
            <input type="checkbox" name="dateFormat" v-model="dateFormat" />
            {{ t("settings.setDateFormat") }}
          </p>
          <h3>{{ t("settings.language") }}</h3>
          <languages
            class="input input--block"
            v-model:locale="locale"
          ></languages>

          <h3>{{ t("settings.aceEditorTheme") }}</h3>
          <AceEditorTheme
            class="input input--block"
            v-model:aceEditorTheme="aceEditorTheme"
            id="aceTheme"
          ></AceEditorTheme>

          <h3>{{ t("settings.autoRename") }}</h3>
          <p>
            <input type="checkbox" name="autoRename" v-model="autoRename" />
            {{ t("settings.autoRenameHelp") }}
          </p>
          <div v-if="autoRename" style="margin-top: 1em;">
            <p>
              <label for="renamePattern">{{ t("settings.renamePattern") }}</label>
              <input
                class="input input--block"
                type="text"
                id="renamePattern"
                v-model="renamePattern"
                :placeholder="t('settings.renamePatternExample')"
              />
            </p>
            <p class="small">{{ t("settings.renamePatternHelp") }}</p>
            <table style="width: 100%; margin: 0.5em 0; font-size: 0.9em; border-collapse: collapse;">
              <thead>
                <tr style="border-bottom: 1px solid var(--divider);">
                  <th style="text-align: left; padding: 4px 8px;">{{ t("settings.placeholder") }}</th>
                  <th style="text-align: left; padding: 4px 8px;">{{ t("settings.description") }}</th>
                  <th style="text-align: left; padding: 4px 8px;">{{ t("settings.example") }}</th>
                </tr>
              </thead>
              <tbody>
                <tr><td style="padding: 3px 8px;"><code>{timestamp}</code></td><td>{{ t("settings.placeholderTimestamp") }}</td><td>20260618153045</td></tr>
                <tr><td style="padding: 3px 8px;"><code>{date}</code></td><td>{{ t("settings.placeholderDate") }}</td><td>20260618</td></tr>
                <tr><td style="padding: 3px 8px;"><code>{time}</code></td><td>{{ t("settings.placeholderTime") }}</td><td>153045</td></tr>
                <tr><td style="padding: 3px 8px;"><code>{year}</code></td><td>{{ t("settings.placeholderYear") }}</td><td>2026</td></tr>
                <tr><td style="padding: 3px 8px;"><code>{month}</code></td><td>{{ t("settings.placeholderMonth") }}</td><td>06</td></tr>
                <tr><td style="padding: 3px 8px;"><code>{day}</code></td><td>{{ t("settings.placeholderDay") }}</td><td>18</td></tr>
                <tr><td style="padding: 3px 8px;"><code>{hour}</code></td><td>{{ t("settings.placeholderHour") }}</td><td>15</td></tr>
                <tr><td style="padding: 3px 8px;"><code>{minute}</code></td><td>{{ t("settings.placeholderMinute") }}</td><td>30</td></tr>
                <tr><td style="padding: 3px 8px;"><code>{second}</code></td><td>{{ t("settings.placeholderSecond") }}</td><td>45</td></tr>
                <tr><td style="padding: 3px 8px;"><code>{name}</code></td><td>{{ t("settings.placeholderName") }}</td><td>文档</td></tr>
                <tr><td style="padding: 3px 8px;"><code>{ext}</code></td><td>{{ t("settings.placeholderExt") }}</td><td>.pdf</td></tr>
                <tr><td style="padding: 3px 8px;"><code>{n}</code></td><td>{{ t("settings.placeholderN") }}</td><td>1, 2, 3</td></tr>
              </tbody>
            </table>
            <p class="small" style="margin-top: 0.5em;">{{ t("settings.renamePatternExamples") }}</p>
            <ul style="font-size: 0.9em; margin: 0.25em 0 0 1.5em;">
              <li><code>{timestamp}{ext}</code> → 20260618153045.pdf</li>
              <li><code>{name}_{timestamp}{ext}</code> → 文档_20260618153045.pdf</li>
              <li><code>{date}_{n}{ext}</code> → 20260618_1.pdf</li>
            </ul>
          </div>
        </div>

        <div class="card-action">
          <input
            class="button button--flat"
            type="submit"
            name="submitProfile"
            :value="t('buttons.update')"
          />
        </div>
      </form>
    </div>

    <div v-if="!noAuth" class="column">
      <form
        class="card"
        v-if="!authStore.user?.lockPassword"
        @submit="updatePassword"
      >
        <div class="card-title">
          <h2>{{ t("settings.changePassword") }}</h2>
        </div>

        <div class="card-content">
          <input
            :class="passwordClass"
            type="password"
            :placeholder="t('settings.newPassword')"
            v-model="password"
            name="password"
          />
          <input
            :class="passwordClass"
            type="password"
            :placeholder="t('settings.newPasswordConfirm')"
            v-model="passwordConf"
            name="passwordConf"
          />
          <input
            v-if="isCurrentPasswordRequired"
            :class="passwordClass"
            type="password"
            :placeholder="t('settings.currentPassword')"
            v-model="currentPassword"
            name="current_password"
            autocomplete="current-password"
          />
        </div>

        <div class="card-action">
          <input
            class="button button--flat"
            type="submit"
            name="submitPassword"
            :value="t('buttons.update')"
          />
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "@/stores/auth";
import { useLayoutStore } from "@/stores/layout";
import { users as api } from "@/api";
import AceEditorTheme from "@/components/settings/AceEditorTheme.vue";
import Languages from "@/components/settings/Languages.vue";
import { computed, inject, onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";
import { authMethod, noAuth } from "@/utils/constants";

const layoutStore = useLayoutStore();
const authStore = useAuthStore();
const { t } = useI18n();

const $showSuccess = inject<IToastSuccess>("$showSuccess")!;
const $showError = inject<IToastError>("$showError")!;

const password = ref<string>("");
const passwordConf = ref<string>("");
const currentPassword = ref<string>("");
const isCurrentPasswordRequired = ref<boolean>(false);
const hideDotfiles = ref<boolean>(false);
const singleClick = ref<boolean>(false);
const redirectAfterCopyMove = ref<boolean>(false);
const dateFormat = ref<boolean>(false);
const locale = ref<string>("");
const aceEditorTheme = ref<string>("");
const autoRename = ref<boolean>(false);
const renamePattern = ref<string>("");

const passwordClass = computed(() => {
  const baseClass = "input input--block";

  if (password.value === "" && passwordConf.value === "") {
    return baseClass;
  }

  if (password.value === passwordConf.value) {
    return `${baseClass} input--green`;
  }

  return `${baseClass} input--red`;
});

onMounted(async () => {
  layoutStore.loading = true;
  if (authStore.user === null) return false;
  locale.value = authStore.user.locale;
  hideDotfiles.value = authStore.user.hideDotfiles;
  singleClick.value = authStore.user.singleClick;
  redirectAfterCopyMove.value = authStore.user.redirectAfterCopyMove;
  dateFormat.value = authStore.user.dateFormat;
  aceEditorTheme.value = authStore.user.aceEditorTheme;
  autoRename.value = authStore.user.autoRename ?? false;
  renamePattern.value = authStore.user.renamePattern ?? "";
  layoutStore.loading = false;
  isCurrentPasswordRequired.value = authMethod == "json";

  return true;
});

const updatePassword = async (event: Event) => {
  event.preventDefault();

  if (
    password.value !== passwordConf.value ||
    password.value === "" ||
    currentPassword.value === "" ||
    authStore.user === null
  ) {
    return;
  }

  try {
    const data = {
      ...authStore.user,
      id: authStore.user.id,
      password: password.value,
    };
    await api.update(data, ["password"], currentPassword.value);
    authStore.updateUser(data);
    $showSuccess(t("settings.passwordUpdated"));
  } catch (e: any) {
    $showError(e);
  } finally {
    password.value = passwordConf.value = "";
  }
};
const updateSettings = async (event: Event) => {
  event.preventDefault();

  try {
    if (authStore.user === null) throw new Error("User is not set!");

    const data = {
      ...authStore.user,
      id: authStore.user.id,
      locale: locale.value,
      hideDotfiles: hideDotfiles.value,
      singleClick: singleClick.value,
      redirectAfterCopyMove: redirectAfterCopyMove.value,
      dateFormat: dateFormat.value,
      aceEditorTheme: aceEditorTheme.value,
      autoRename: autoRename.value,
      renamePattern: renamePattern.value,
    };

    await api.update(data, [
      "locale",
      "hideDotfiles",
      "singleClick",
      "autoRename",
      "renamePattern",
      "redirectAfterCopyMove",
      "dateFormat",
      "aceEditorTheme",
    ]);
    authStore.updateUser(data);
    $showSuccess(t("settings.settingsUpdated"));
  } catch (err) {
    if (err instanceof Error) {
      $showError(err);
    }
  }
};
</script>
