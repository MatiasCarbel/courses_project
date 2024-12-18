import { formatCookies, getCookieValue } from "@/lib/api.utils";
import { jwtDecode } from "jwt-decode";
import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  const baseUrl = process.env.NEXT_PUBLIC_COURSES_API_URL ?? "";

  try {
    const cookie = formatCookies();
    const cookieValue = getCookieValue();

    if (!cookieValue) {
      return NextResponse.json(
        { message: "Not authenticated", courses: [] },
        { status: 401 }
      );
    }

    const decoded = jwtDecode(cookieValue ?? "") as any;
    const userId = decoded?.user_id;

    if (!userId) {
      return NextResponse.json(
        { message: "Invalid user ID", courses: [] },
        { status: 400 }
      );
    }

    const url = `${baseUrl}/courses/myCourses`;

    const coursesReq = await fetch(url, {
      credentials: "include",
      headers: {
        Cookie: cookie,
        Authorization: `Bearer ${cookieValue}`,
      },
    });

    console.log("coursesReq: ", coursesReq);

    if (!coursesReq.ok) {
      return NextResponse.json(
        {
          message: `Error fetching courses: ${coursesReq.status} ${coursesReq.statusText}`,
          courses: [],
        },
        { status: coursesReq.status }
      );
    }

    const coursesJson = await coursesReq.json();
    console.log("coursesJson: ", coursesJson);

    return NextResponse.json(
      {
        message: "Courses fetched",
        data: {
          courses: coursesJson?.data?.courses ?? [],
        },
      },
      { status: 200 }
    );
  } catch (error) {
    console.error("Error in myCourses route:", error);
    return NextResponse.json(
      { message: "Internal server error", courses: [] },
      { status: 500 }
    );
  }
}
