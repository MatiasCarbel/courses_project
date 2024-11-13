import { formatCookies, getCookieValue } from "@/lib/api.utils";
import { NextRequest, NextResponse } from "next/server";
import { jwtDecode } from "jwt-decode";

export async function POST(request: NextRequest) {
  const baseUrl = process.env.NEXT_PUBLIC_USERS_API_URL ?? "";

  const data = await request.json();
  const { courseId, comment } = data;

  // Get courseId from query params.
  const cookie = formatCookies();
  const cookieValue = getCookieValue();

  const decoded = jwtDecode(cookieValue ?? "") as any;
  const userId = decoded?.id;

  const commentReq = await fetch(`${baseUrl}/comments`, {
    method: "POST",
    headers: {
      cookie: cookie,
    },
    body: JSON.stringify({
      course_id: Number(courseId),
      user_id: Number(userId),
      comment: comment,
    }),
  });
  const commentJson = await commentReq.json();

  return NextResponse.json(
    { message: "comment added", comment: commentJson },
    { status: 200 }
  );
}
