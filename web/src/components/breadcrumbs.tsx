import { A } from "@solidjs/router";
import { BaseInfo } from "@sv/apis/media/models";
import {
  CLIENT_URLS,
  ROOT_FOLDER_ID,
  ROOT_FOLDER_NAME,
} from "@sv/utils/consts";
import { For, Show } from "solid-js";

interface Props {
  ancestors: BaseInfo[];
  currentFolder: BaseInfo;
}

function Breadcrumbs(props: Props) {
  return (
    <div class="text-primary">
      <Show when={props.currentFolder.id !== ROOT_FOLDER_ID}>
        <span>
          <A href={CLIENT_URLS.DRIVE} class="link">
            {ROOT_FOLDER_NAME}
          </A>
          {" / "}
        </span>
      </Show>
      <For each={props.ancestors.toReversed()}>
        {(ancestor) => (
          <span>
            <A href={CLIENT_URLS.DRIVE + "/" + ancestor.id} class="link">
              {ancestor.name}
            </A>
            {" / "}
          </span>
        )}
      </For>
      <span class="font-bold text-neutral">{props.currentFolder.name}</span>
    </div>
  );
}

export default Breadcrumbs;
