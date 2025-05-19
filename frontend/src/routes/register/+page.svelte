<script lang="ts">
  import { goto } from "$app/navigation";
  import { ROLE_ADMIN } from "$lib/constants/constants";
  import { register } from "$lib/services/auth";
  import type { UserJWT } from "$lib/types/user";
  import { jwtDecode } from "jwt-decode";
  import { onMount } from "svelte";

  let first_name = "";
  let last_name = "";
  let email = "";
  let password = "";
  let error = "";

  async function handleRegister(event: Event) {
    event.preventDefault();
    if (!first_name || !last_name || !email || !password) {
      error = "Please fill in all fields.";
      return;
    }
    error = "";
    const data = await register(first_name, last_name, email, password);
    if (data.success === "success") {
      // Optionally, auto-login after registration:
      localStorage.setItem("token", data.data);
      const decoded: UserJWT = jwtDecode(data.data);
      if (decoded.role === ROLE_ADMIN) {
        goto("/admin/dashboard");
      } else {
        goto("/employee/dashboard");
      }
    } else {
      error = data.message || "Registration failed";
    }
  }

  onMount(() => {
    const token = localStorage.getItem("token") as string;

    if (token) {
      const decoded: UserJWT = jwtDecode(token);

      if (decoded.role === ROLE_ADMIN) {
        goto("/admin/dashboard");
        return;
      }

      goto("/employee/dashboard");
    }
  });
</script>

<div class="flex items-center justify-center min-h-screen bg-gray-50">
  <div class="w-full max-w-md p-8 space-y-6 bg-white rounded shadow">
    <h2 class="text-2xl font-bold text-center text-gray-900">
      Make your account
    </h2>
    {#if error}
      <div class="p-2 mb-2 text-sm text-red-700 bg-red-100 rounded">
        {error}
      </div>
    {/if}
    <form class="space-y-4" on:submit|preventDefault={handleRegister}>
      <div>
        <label
          for="first_name"
          class="block mb-2 text-sm font-medium text-gray-700">First Name</label
        >
        <input
          id="first_name"
          type="text"
          bind:value={first_name}
          required
          class="block w-full px-3 py-2 border border-gray-300 rounded focus:ring-primary-500 focus:border-primary-500"
          placeholder="First Name"
        />
      </div>
      <div>
        <label
          for="last_name"
          class="block mb-2 text-sm font-medium text-gray-700">Last Name</label
        >
        <input
          id="last_name"
          type="text"
          bind:value={last_name}
          required
          class="block w-full px-3 py-2 border border-gray-300 rounded focus:ring-primary-500 focus:border-primary-500"
          placeholder="Last Name"
        />
      </div>
      <div>
        <label for="email" class="block mb-2 text-sm font-medium text-gray-700"
          >Email</label
        >
        <input
          id="email"
          type="email"
          bind:value={email}
          required
          class="block w-full px-3 py-2 border border-gray-300 rounded focus:ring-primary-500 focus:border-primary-500"
          placeholder="name@company.com"
        />
      </div>
      <div>
        <label
          for="password"
          class="block mb-2 text-sm font-medium text-gray-700">Password</label
        >
        <input
          id="password"
          type="password"
          bind:value={password}
          required
          class="block w-full px-3 py-2 border border-gray-300 rounded focus:ring-primary-500 focus:border-primary-500"
          placeholder="••••••••"
        />
      </div>
      <button
        type="submit"
        class="w-full px-4 py-2 text-white bg-blue-600 rounded hover:bg-blue-700 focus:ring-4 focus:ring-blue-300 font-medium"
      >
        Register
      </button>
      <p class="text-sm text-center text-gray-500">
        Already have an account?
        <a href="/login" class="text-blue-600 hover:underline">Login</a>
      </p>
    </form>
  </div>
</div>
