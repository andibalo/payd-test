<script lang="ts">
  import {
    ChartOutline,
    GridSolid,
    MailBoxSolid,
    UserSolid,
  } from "flowbite-svelte-icons";
  import { SidebarMenu } from "$lib";
  import type { SvelteComponent } from "svelte";
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { jwtDecode } from "jwt-decode";
  import type { UserJWT } from "$lib/types/user";
  import { ROLE_ADMIN } from "$lib/constants/constants";

  const adminSidebarOptions = [
    {
      label: "Dashboard",
      href: "/admin/dashboard",
      icon: ChartOutline as unknown as typeof SvelteComponent,
    },
    {
      label: "Shifts",
      href: "/admin/dashboard/shifts",
      icon: GridSolid as unknown as typeof SvelteComponent,
    },
    {
      label: "Shift Requests",
      href: "/admin/dashboard/shift-requests",
      icon: MailBoxSolid as unknown as typeof SvelteComponent,
    },
    {
      label: "Shift Assignments",
      href: "/admin/dashboard/shift-assignments",
      icon: UserSolid as unknown as typeof SvelteComponent,
    },
  ];

  onMount(() => {
    const token = localStorage.getItem("token") as string;

    if (!token) {
      goto("/login");
      return;
    }

    const decoded: UserJWT = jwtDecode(token);

    if (decoded.role !== ROLE_ADMIN) {
      goto("/login");
    }
  });
</script>

<div class="flex min-h-screen">
  <SidebarMenu options={adminSidebarOptions} sidebarTitle="Admin Dashboard" />
  <main class="flex-1 ml-64 h-screen overflow-y-auto p-8 bg-gray-50">
    <slot />
  </main>
</div>
