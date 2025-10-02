import { createContext } from "solid-js";

export interface SelectedItem {
  id: string;
  type: string;
  name: string;
}

export interface CtxType {
  handleTap: (tappedItem: SelectedItem, singleTapAction?: () => void) => void;
  selectedItem: () => SelectedItem | null;
  clearSelection: () => void;
}

const CTX = createContext<CtxType>();

export default CTX;
