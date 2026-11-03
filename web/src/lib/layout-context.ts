import { useOutletContext } from "react-router-dom";

export type AppOutletContext = { searchQuery: string };

export function useAppOutletContext() {
  return useOutletContext<AppOutletContext>();
}
