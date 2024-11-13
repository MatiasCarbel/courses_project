import { formatCookies } from "@/lib/api.utils";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const baseUrl = process.env.NEXT_PUBLIC_USERS_API_URL ?? "";

  const formData = await request.formData();
  const courseId = formData.get("courseId");

  const cookie = formatCookies();

  const resourceReq = await fetch(`${baseUrl}/upload/${courseId}`, {
    method: "POST",
    headers: {
      cookie: cookie,
    },
    body: formData,
  });

  const resourceText = await resourceReq.text();

  if (!resourceReq.ok) {
    return NextResponse.json(
      { message: "Error while uploading resource" },
      { status: 400 }
    );
  }

  let resourceJson;
  try {
    resourceJson = JSON.parse(resourceText);
  } catch (error) {
    return NextResponse.json(
      { message: "Error while parsing server response" },
      { status: 500 }
    );
  }

  return NextResponse.json(
    { message: "Resource Added", course: resourceJson },
    { status: 200 }
  );
}
