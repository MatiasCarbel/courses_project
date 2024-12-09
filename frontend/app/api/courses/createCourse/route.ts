import { formatCookies, getCookieValue } from "@/lib/api.utils";
import { NextRequest, NextResponse } from "next/server";
import { jwtDecode } from "jwt-decode";

export async function POST(request: NextRequest) {
  const baseUrl = process.env.NEXT_PUBLIC_USERS_API_URL ?? "";

  const data = await request.json();
  const {
    courseImage,
    courseName,
    courseDescription,
    courseDuration,
    courseCategory,
    courseRequirements,
  } = data;

  // Get courseId from query params.
  const cookie = formatCookies();
  const cookieValue = getCookieValue();

  const decoded = jwtDecode(cookieValue ?? "") as any;
  const userId = decoded?.id;

  const body = {
    title: courseName,
    description: courseDescription,
    instructor: decoded?.username || "Unknown Instructor",
    category: courseCategory,
    duration: Number(courseDuration),
    available_seats: 20, // Default value or make it configurable
    image_url: courseImage,
  };

  const courseReq = await fetch(`${baseUrl}/courses`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      cookie: cookie,
    },
    body: JSON.stringify(body),
  });
  const courseJson = await courseReq.json();

  if (courseJson.error) {
    return NextResponse.json(
      { message: "Error creating course", error: courseJson.error },
      { status: 400 }
    );
  }

  return NextResponse.json(
    { message: "Course created", course: courseJson },
    { status: 200 }
  );
}
