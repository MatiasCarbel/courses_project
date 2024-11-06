import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { cookies } from "next/headers";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatCookies() {
  const cookiesStore = cookies();
  const cookie = cookiesStore.get("auth");

  const cookieValue = cookie?.value;
  if (!cookieValue) return ``;

  return `auth=${cookieValue}`;
}

export function getCookieValue() {
  const cookiesStore = cookies();
  const cookie = cookiesStore.get("auth");

  return cookie?.value;
}
