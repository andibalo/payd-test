<script lang="ts">
  import { Sidebar, SidebarGroup, SidebarItem } from "flowbite-svelte";
  import type { SvelteComponent } from "svelte";
  import { page } from "$app/stores";
  import { get } from "svelte/store";

  export let options: {
    label: string;
    href: string;
    icon?: typeof SvelteComponent;
    subtext?: string;
    subtextClass?: string;
  }[] = [];

  export let isOpen: boolean = true;
  export let closeSidebar: () => void = () => {};
  export let sidebarTitle: string = "Dashboard";

  $: activeUrl = get(page).url.pathname;
</script>

<Sidebar
  {activeUrl}
  backdrop={false}
  params={{ x: -50, duration: 50 }}
  class="z-50 h-full"
  position="absolute"
  activeClass="p-2"
  nonActiveClass="p-2"
>
  <SidebarGroup>
    <div class="p-6 text-2xl font-bold border-b border-gray-700">
      {sidebarTitle}
    </div>
    {#each options as option}
      <SidebarItem label={option.label} href={option.href}>
        {#if option.icon}
          <svelte:component
            this={option.icon}
            class="h-5 w-5 text-gray-500 transition duration-75 group-hover:text-gray-900 dark:text-gray-400 dark:group-hover:text-white"
          />
        {/if}
        {#if option.subtext}
          <span
            class={option.subtextClass ??
              "ms-3 inline-flex items-center justify-center rounded-full bg-gray-200 px-2 text-sm font-medium text-gray-800 dark:bg-gray-700 dark:text-gray-300"}
          >
            {option.subtext}
          </span>
        {/if}
      </SidebarItem>
    {/each}
  </SidebarGroup>
</Sidebar>
