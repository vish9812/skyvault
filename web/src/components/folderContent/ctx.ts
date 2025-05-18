import { createContext } from "solid-js";

interface CtxType {
  handleTap: (type: string, id: string, singleTapAction?: () => void) => void;
}

const CTX = createContext<CtxType>();

export default CTX;
