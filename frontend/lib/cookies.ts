import { cookies } from "next/headers";

export const formatCookies = () => {
  const cookieStore = cookies();
  const cookieArray = cookieStore.getAll();
  return cookieArray
    .map((cookie) => `${cookie.name}=${cookie.value}`)
    .join("; ");
};

export const getCookieValue = () => {
  const cookieStore = cookies();
  const token = cookieStore.get("token");
  return token?.value;
};
