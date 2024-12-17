import { formatCookies, getCookieValue } from "@/lib/api.utils";
import { jwtDecode } from "jwt-decode";
import { NextRequest, NextResponse } from "next/server";

export async function DELETE(
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

    const response = await fetch(`${baseUrl}/courses/${courseId}`, {
      method: "DELETE",
      headers: {
        cookie: cookie,
      },
    });

    if (!response.ok) {
      throw new Error("Failed to delete course");
    }

    return NextResponse.json(
      { message: "Course deleted successfully" },
      { status: 200 }
    );
  } catch (error: any) {
    console.error("Error deleting course:", error);
    return NextResponse.json(
      { message: "Error deleting course", error: error.message },
      { status: 500 }
    );
  }
}
