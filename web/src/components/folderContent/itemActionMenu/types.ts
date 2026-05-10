import type { FileInfo, FolderInfo } from "@sv/apis/media/models";
import { FOLDER_CONTENT_TYPES } from "@sv/utils/consts";

export type ItemType =
  | typeof FOLDER_CONTENT_TYPES.FILE
  | typeof FOLDER_CONTENT_TYPES.FOLDER;

export type ActionItem = FileInfo & FolderInfo;

export interface DialogProps {
  type: ItemType;
  item: ActionItem;
  itemLabel: string;
  open: boolean;
  onClose: () => void;
  onDone: () => void;
}

