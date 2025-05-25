import { createContext } from "solid-js";

export interface CtxType {
  handleTap: (type: string, id: string, singleTapAction?: () => void) => void;
}

const CTX = createContext<CtxType>();

export default CTX;
