import { formatCookies, getCookieValue } from "@/lib/api.utils";
import { NextRequest, NextResponse } from "next/server";
import { jwtDecode } from "jwt-decode";

export async function POST(request: NextRequest) {
  const baseUrl =
    process.env.NEXT_PUBLIC_COURSES_API_URL ?? "http://localhost:8002";

  try {
    const data = await request.json();
    const { courseId } = data;

    // Get courseId from query params.
    const cookie = formatCookies();
    const cookieValue = getCookieValue();

    if (!cookieValue) {
      return NextResponse.json(
        { message: "Authentication required" },
        { status: 401 }
      );
    }

    const decoded = jwtDecode(cookieValue ?? "") as any;
    const userId = decoded?.user_id;

    if (!courseId || !userId) {
      return NextResponse.json(
        { message: "Course ID and User ID are required" },
        { status: 400 }
      );
    }

    const courseReq = await fetch(`${baseUrl}/enrollments`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        cookie: cookie,
      },
      body: JSON.stringify({
        course_id: courseId,
        user_id: userId,
      }),
    });

    const courseJson = await courseReq.json();

    // Handle error responses from the backend
    if (!courseReq.ok) {
      return NextResponse.json(
        { message: courseJson.error || "Failed to enroll in course" },
        { status: courseReq.status }
      );
    }

    return NextResponse.json(
      { message: "Successfully enrolled in course", course: courseJson },
      { status: 200 }
    );
  } catch (error: any) {
    console.error("Error in course enrollment:", error);
    return NextResponse.json(
      { message: "Internal server error" },
      { status: 500 }
    );
  }
}
