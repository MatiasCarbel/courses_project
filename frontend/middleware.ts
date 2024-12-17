import { NextRequest, NextResponse } from "next/server";
import { jwtDecode, JwtPayload } from "jwt-decode";
import routes from "./lib/routes";

export async function middleware(request: NextRequest) {
  const url = request.nextUrl.clone();
  const cookie = request.cookies.get("token");
  const cookieValue = cookie?.value ?? "";
  let isAuthenticated = false;

  if (cookieValue !== "") {
    try {
      const decoded = jwtDecode<JwtPayload>(cookieValue);

      console.log("decoded: ", decoded);
      console.log("decoded.exp: ", decoded.exp);
      console.log("Date.now(): ", Date.now());

      // JWT exp is in seconds, Date.now() is in milliseconds
      if (decoded.exp && decoded.exp * 1000 > Date.now()) {
        isAuthenticated = true;
      }
    } catch (error) {
      console.error("Error decoding token:", error);
      isAuthenticated = false;
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
  const courseRouteRegex = new RegExp(
    `^${routes.course.replace(":id", "[^/]+")}$`
  );

  console.log(
    "courseRouteRegex.test(url.pathname): ",
    courseRouteRegex.test(url.pathname)
  );

  console.log("isAuthenticated: ", isAuthenticated);

  if (
    !isAuthenticated &&
    [routes.myCourses, routes.course].some(
      (route) => url.pathname === route || courseRouteRegex.test(url.pathname)
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
