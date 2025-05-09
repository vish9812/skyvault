import { Meter } from "@kobalte/core/meter";
import type { JSX } from "solid-js";

interface StorageMeterProps {
  used: number;
  total: number;
}

const StorageMeter = (props: StorageMeterProps): JSX.Element => {
  return (
    <div class="flex-1">
      <div class="flex justify-between mb-1">
        <span class="text-sm font-medium text-gray-700">Storage Usage</span>
        <span class="text-sm text-gray-500">
          {props.used} GB of {props.total} GB used
        </span>
      </div>
      <Meter
        value={props.used}
        minValue={0}
        maxValue={props.total}
        class="w-full"
      >
        <Meter.Track class="h-2 bg-gray-200 rounded-full overflow-hidden">
          <Meter.Fill class="bg-gradient-to-r from-indigo-500 to-sky-400 h-full rounded-full transition-all" />
        </Meter.Track>
      </Meter>
    </div>
  );
};

export default StorageMeter;
