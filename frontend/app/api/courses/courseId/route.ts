import { formatCookies } from "@/lib/api.utils";
import { NextRequest, NextResponse } from "next/server";
import { cookies } from "next/headers";
import { jwtDecode } from "jwt-decode";

export async function GET(request: NextRequest) {
  const cookiesStore = cookies();
  const authCookie = cookiesStore.get("auth");

  const cookieValue = authCookie?.value;
  if (!cookieValue)
    return NextResponse.json(
      { message: "Not authenticated.", shouldLogin: true },
      { status: 200 }
    );

  const decoded = jwtDecode(cookieValue ?? "") as any;
  const userId = decoded?.id;

  const baseUrl = process.env.NEXT_PUBLIC_BASE_API_URL ?? "";

  // Get courseId from query params.
  const cookie = formatCookies();
  const courseId = request.nextUrl.searchParams.get("courseId");
  const courseReq = await fetch(
    `${baseUrl}/course/${courseId}?userId=${userId}`,
    {
      headers: {
        cookie: cookie,
      },
    }
  );
  const courseJson = await courseReq.json();

  const courseCommentsReq = await fetch(`${baseUrl}/comments/${courseId}`);
  const courseCommentsJson = await courseCommentsReq.json();

  console.log(courseCommentsJson);

  return NextResponse.json(
    {
      message: "Courses fetched",
      course: courseJson,
      comments: courseCommentsJson,
    },
    { status: 200 }
  );
}
