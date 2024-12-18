import { formatCookies, getCookieValue } from "@/lib/api.utils";
import { jwtDecode } from "jwt-decode";
import { NextRequest, NextResponse } from "next/server";

export async function GET(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  const baseUrl = process.env.NEXT_PUBLIC_COURSES_API_URL ?? "";
  const courseId = params.id;

  try {
    const cookie = formatCookies();

    // Get course details
    const courseResponse = await fetch(`${baseUrl}/courses/${courseId}`);
    if (!courseResponse.ok) {
      throw new Error("Failed to fetch course");
    }
    const courseData = await courseResponse.json();

    // Check enrollment if user is authenticated
    if (cookie) {
      const enrollmentResponse = await fetch(
        `${baseUrl}/enrollments/check/${courseId}`,
        {
          headers: {
            Authorization: `Bearer ${cookie.split("=")[1]}`,
          },
        }
      );

      if (enrollmentResponse.ok) {
        const enrollmentData = await enrollmentResponse.json();
        return NextResponse.json({
          course: {
            ...courseData.data,
            is_subscribed: enrollmentData?.data?.enrolled ?? false,
          },
        });
      }
    }

    return NextResponse.json({
      course: {
        ...courseData.data,
        is_subscribed: false,
      },
    });
  } catch (error: any) {
    return NextResponse.json(
      { message: error.message || "Error fetching course" },
      { status: 500 }
    );
  }
}

export async function PUT(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  const baseUrl = process.env.NEXT_PUBLIC_COURSES_API_URL ?? "";
  const cookieValue = getCookieValue();

  if (!cookieValue) {
    return NextResponse.json({ message: "Not authenticated" }, { status: 401 });
  }

  const decoded = jwtDecode(cookieValue) as any;
  if (!decoded?.admin) {
    return NextResponse.json({ message: "Not authorized" }, { status: 403 });
  }

  try {
    const body = await request.json();

    // Ensure numeric fields are numbers
    const updatedBody = {
      ...body,
      duration: Number(body.duration),
      available_seats: Number(body.available_seats),
    };

    const response = await fetch(`${baseUrl}/courses/${params.id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${cookieValue}`,
      },
      body: JSON.stringify(updatedBody),
    });

    const data = await response.json();

    if (!response.ok) {
      return NextResponse.json(
        { message: data.error || "Failed to update course" },
        { status: response.status }
      );
    }

    return NextResponse.json({
      message: "Course updated successfully",
      course: data,
    });
  } catch (error: any) {
    console.error("Error updating course:", error);
    return NextResponse.json(
      { message: error.message || "Error updating course" },
      { status: 500 }
    );
  }
}

export async function DELETE(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  const baseUrl = process.env.NEXT_PUBLIC_COURSES_API_URL ?? "";
  const cookieValue = getCookieValue();

  if (!cookieValue) {
    return NextResponse.json({ message: "Not authenticated" }, { status: 401 });
  }

  const decoded = jwtDecode(cookieValue) as any;
  if (!decoded?.admin) {
    return NextResponse.json({ message: "Not authorized" }, { status: 403 });
  }

  try {
    const response = await fetch(`${baseUrl}/courses/${params.id}`, {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${cookieValue}`,
      },
    });

    if (!response.ok) {
      throw new Error("Failed to delete course");
    }

    return NextResponse.json({ message: "Course deleted successfully" });
  } catch (error: any) {
    return NextResponse.json(
      { message: error.message || "Error deleting course" },
      { status: 500 }
    );
  }
}
