import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const baseUrl =
    process.env.NEXT_PUBLIC_COURSES_API_URL ?? "http://courses-api:8002";

  try {
    const courseIds = await request.json();

    const availabilityRes = await fetch(`${baseUrl}/courses/availability`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(courseIds),
    });

    if (!availabilityRes.ok) {
      throw new Error("Failed to fetch availabilities");
    }

    const availabilityData = await availabilityRes.json();
    return NextResponse.json(availabilityData, { status: 200 });
  } catch (error: any) {
    console.error("Error checking availability:", error);
    return NextResponse.json(
      { message: "Error checking availability", error: error.message },
      { status: 500 }
    );
  }
}
