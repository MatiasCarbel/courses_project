import { formatCookies } from "@/lib/api.utils";
import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  console.log("fetching course details");
  const coursesApiUrl =
    process.env.NEXT_PUBLIC_COURSES_API_URL ?? "http://courses-api:8002";

  const courseId = request.nextUrl.searchParams.get("courseId");

  try {
    const searchResponse = await fetch(`${coursesApiUrl}/courses/${courseId}`, {
      headers: {
        Accept: "application/json",
      },
      cache: "no-cache",
    });

    const course = await searchResponse.json();
    console.log("course: ", course);

    if (!course) {
      return NextResponse.json(
        { message: "Course not found" },
        { status: 404 }
      );
    }

    console.log("fetching enrollment data");

    const cookie = formatCookies();
    if (!cookie) {
      return NextResponse.json(
        {
          message: "Course details fetched",
          course: {
            ...course,
            is_subscribed: false,
          },
        },
        { status: 200 }
      );
    }

    const enrollmentResponse = await fetch(
      `${coursesApiUrl}/enrollments/check/${courseId}`,
      {
        headers: {
          Authorization: `Bearer ${cookie.split("=")[1]}`,
          Accept: "application/json",
        },
      }
    );

    if (!enrollmentResponse.ok) {
      console.error("Enrollment check failed:", enrollmentResponse.status);
      return NextResponse.json(
        {
          message: "Course details fetched",
          course: {
            ...course,
            is_subscribed: false,
          },
        },
        { status: 200 }
      );
    }

    const enrollmentData = await enrollmentResponse.json();
    console.log("enrollmentData: ", enrollmentData);

    return NextResponse.json(
      {
        message: "Course details fetched",
        course: {
          ...course,
          is_subscribed: enrollmentData?.data?.enrolled ?? false,
        },
      },
      { status: 200 }
    );
  } catch (error: any) {
    console.error("Error fetching course details:", error);
    return NextResponse.json(
      { message: "Error fetching course details", error: error.message },
      { status: 500 }
    );
  }
}
