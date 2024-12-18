import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { cookies } from "next/headers";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const formatCookies = () => {
  if (typeof document === "undefined") {
    // Server-side
    const cookies = require("next/headers").cookies;
    const cookieStore = cookies();
    const token = cookieStore.get("token");
    return token ? `token=${token.value}` : "";
  }
  // Client-side
  return document.cookie;
};

export function getCookieValue() {
  const cookiesStore = cookies();
  const cookie = cookiesStore.get("token");

  return cookie?.value;
}
