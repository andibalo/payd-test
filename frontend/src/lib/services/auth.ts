import { PUBLIC_RMS_API_BASE_URL } from '$env/static/public';

export async function login(email: string, password: string) {
  const res = await fetch(`${PUBLIC_RMS_API_BASE_URL}/api/v1/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });
  const data = await res.json();
  if (data.success === 'success' && data.data) {
    localStorage.setItem('token', data.data);
  }
  return data;
}

export async function register(first_name: string, last_name: string, email: string, password: string) {
  const res = await fetch(`${PUBLIC_RMS_API_BASE_URL}/api/v1/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ first_name, last_name, email, password })
  });

  const data = await res.json();
  if (data.success === 'success' && data.data) {
    localStorage.setItem('token', data.data);
  }
  return data;
}