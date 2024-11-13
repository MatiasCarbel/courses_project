import { formatCookies, getCookieValue } from "@/lib/api.utils";
import { NextRequest, NextResponse } from "next/server";
import { jwtDecode } from "jwt-decode";

export async function POST(request: NextRequest) {
  const baseUrl = process.env.NEXT_PUBLIC_USERS_API_URL ?? "";

  const data = await request.json();
  const { courseId } = data;

  // Get courseId from query params.
  const cookie = formatCookies();
  const cookieValue = getCookieValue();

  const decoded = jwtDecode(cookieValue ?? "") as any;
  const userId = decoded?.id;

  const courseReq = await fetch(`${baseUrl}/subscription`, {
    method: "POST",
    headers: {
      cookie: cookie,
    },
    body: JSON.stringify({
      course_id: Number(courseId),
      user_id: Number(userId),
    }),
  });
  const courseJson = await courseReq.json();

  return NextResponse.json(
    { message: "Courses fetched", course: courseJson },
    { status: 200 }
  );
}
