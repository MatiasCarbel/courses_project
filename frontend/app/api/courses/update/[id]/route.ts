import { formatCookies, getCookieValue } from "@/lib/api.utils";
import { jwtDecode } from "jwt-decode";
import { NextRequest, NextResponse } from "next/server";

export async function PUT(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  const baseUrl = process.env.NEXT_PUBLIC_COURSES_API_URL ?? "";
  const courseId = params.id;

  try {
    const cookie = formatCookies();
    const cookieValue = getCookieValue();

    if (!cookieValue) {
      return NextResponse.json(
        { message: "Not authenticated" },
        { status: 401 }
      );
    }

    const decoded = jwtDecode(cookieValue ?? "") as any;
    if (!decoded?.admin) {
      return NextResponse.json({ message: "Not authorized" }, { status: 403 });
    }

    const data = await request.json();
    const { title, description, image_url } = data;

    // Get the current course first to preserve other fields
    const currentCourse = await fetch(`${baseUrl}/courses/${courseId}`, {
      headers: {
        Accept: "application/json",
        cookie: cookie,
      },
    }).then((res) => res.json());

    // Merge the updates with existing data
    const updatedCourse = {
      ...currentCourse,
      title: title || currentCourse.title,
      description: description || currentCourse.description,
      image_url: image_url || currentCourse.image_url,
    };

    const response = await fetch(`${baseUrl}/courses/${courseId}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        cookie: cookie,
      },
      body: JSON.stringify(updatedCourse),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || "Failed to update course");
    }

    const responseData = await response.json();

    return NextResponse.json(
      { message: "Course updated successfully", course: responseData },
      { status: 200 }
    );
  } catch (error: any) {
    console.error("Error updating course:", error);
    return NextResponse.json(
      { message: error.message || "Error updating course" },
      { status: 500 }
    );
  }
}
