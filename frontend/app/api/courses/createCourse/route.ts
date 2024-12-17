import { formatCookies, getCookieValue } from "@/lib/api.utils";
import { NextRequest, NextResponse } from "next/server";
import { jwtDecode } from "jwt-decode";

export async function POST(request: NextRequest) {
  const baseUrl = process.env.NEXT_PUBLIC_COURSES_API_URL ?? "";

  const data = await request.json();
  const {
    courseImage,
    courseName,
    courseDescription,
    courseDuration,
    courseCategory,
    courseInstructor,
    availableSeats,
  } = data;

  const cookie = formatCookies();
  const cookieValue = getCookieValue();

  if (!cookieValue) {
    return NextResponse.json({ message: "Not authenticated" }, { status: 401 });
  }

  const decoded = jwtDecode(cookieValue ?? "") as any;
  if (!decoded?.admin) {
    return NextResponse.json({ message: "Not authorized" }, { status: 403 });
  }

  const body = {
    title: courseName,
    description: courseDescription,
    instructor: courseInstructor,
    category: courseCategory,
    duration: Number(courseDuration),
    available_seats: Number(availableSeats),
    image_url: courseImage,
  };

  try {
    const courseReq = await fetch(`${baseUrl}/courses`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        cookie: cookie,
      },
      body: JSON.stringify(body),
    });

    if (!courseReq.ok) {
      const error = await courseReq.json();
      return NextResponse.json(
        { message: "Error creating course", error: error.error },
        { status: courseReq.status }
      );
    }

    const courseJson = await courseReq.json();
    return NextResponse.json(
      { message: "Course created", course: courseJson },
      { status: 200 }
    );
  } catch (error: any) {
    return NextResponse.json(
      { message: "Error creating course", error: error.message },
      { status: 500 }
    );
  }
}
