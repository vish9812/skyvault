import useAppCtx from "@sv/store/appCtxProvider";
import Format from "@sv/utils/format";
import { Show } from "solid-js";

export default function StorageQuotaBar() {
  const appCtx = useAppCtx();

  const storageUsage = () => appCtx.storageUsage();
  const usagePercentage = () => {
    const usage = storageUsage();
    if (usage.quotaBytes === 0) return 0;
    return Math.min((usage.usedBytes / usage.quotaBytes) * 100, 100);
  };

  const isNearLimit = () => usagePercentage() >= 80;
  const isAtLimit = () => usagePercentage() >= 95;

  return (
    <div class="px-4 py-3 border-t border-border">
      <div class="text-xs text-neutral-light mb-1 flex justify-between items-center">
        <span>Storage</span>
        <span>
          {Format.size(storageUsage().usedBytes)} /{" "}
          {Format.size(storageUsage().quotaBytes)}
        </span>
      </div>

      {/* Progress bar */}
      <div class="w-full bg-bg-muted rounded-full h-2">
        <div
          classList={{
            "h-2 rounded-full transition-all duration-300": true,
            "bg-primary": !isNearLimit(),
            "bg-warning": isNearLimit() && !isAtLimit(),
            "bg-error": isAtLimit(),
          }}
          style={{ width: `${usagePercentage()}%` }}
        />
      </div>

      {/* Warning message */}
      <Show when={isNearLimit()}>
        <div class="text-xs text-warning mt-1">
          <Show when={isAtLimit()} fallback="Storage space running low">
            Storage quota almost full
          </Show>
        </div>
      </Show>
    </div>
  );
}
