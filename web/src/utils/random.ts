import { v7 } from "uuid";

function id() {
  return v7();
}

const Random = {
  id,
} as const;

export default Random;
