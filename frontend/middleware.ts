import { NextRequest, NextResponse } from "next/server";
import { jwtDecode } from "jwt-decode";
import routes from "./lib/routes";

export async function middleware(request: NextRequest) {
  const url = request.nextUrl.clone();
  const cookie = request.cookies.get("auth");
  const cookieValue = cookie?.value ?? "";
  let isAuthenticated = false;

  if (cookieValue !== "") {
    const decoded = jwtDecode(cookieValue ?? "");

    // Check expiration time against current itme.
    if ((decoded?.exp ?? 0) * 1000 < Date.now()) {
      isAuthenticated = false;
    } else {
      isAuthenticated = true;
    }
  }

  // Handle public or static content early
  if (url.pathname.startsWith("/_next/") || url.pathname === "/favicon.ico") {
    return NextResponse.next();
  }

  // Redirect logic for authenticated users trying to access login while logged in
  if (
    isAuthenticated &&
    (url.pathname.includes(routes.login) ||
      url.pathname.includes(routes.register))
  ) {
    const pathToRedirect = url.searchParams.get("redirect");
    const redirectPath = pathToRedirect ?? routes.home;
    return NextResponse.redirect(new URL(redirectPath, request.url));
  }

  if (url.pathname === "/") {
    return NextResponse.redirect(new URL(routes.home, request.url));
  }

  // Auth-required pages for unauthenticated users
  if (
    !isAuthenticated &&
    [routes.myCourses, routes.course].some((route) =>
      url.pathname.includes(route)
    )
  ) {
    return NextResponse.redirect(new URL(routes.login, request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    "/",
    "/register",
    "/courses",
    "/course",
    "/home",
    "/logout",
    "/login",
    "/api/:path*",
    "/((?!_next/static|_next/image|favicon.ico).*)",
  ],
};
