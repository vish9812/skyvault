const API_BASE = "http://localhost:8090/api/v1/pub/auth";

export async function signIn({
  email,
  password,
}: {
  email: string;
  password: string;
}) {
  const res = await fetch(`${API_BASE}/sign-in`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      provider: "email",
      providerUserId: email,
      password,
    }),
  });
  const data = await res.json();
  if (!res.ok) throw data;
  return data;
}

export async function signUp({
  fullName,
  email,
  password,
}: {
  fullName: string;
  email: string;
  password: string;
}) {
  const res = await fetch(`${API_BASE}/sign-up`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      fullName,
      email,
      password,
      provider: "email",
    }),
  });
  const data = await res.json();
  if (!res.ok) throw data;
  return data;
}
