<script lang="ts">
  import {
    Breadcrumb,
    BreadcrumbItem,
    Button,
    ButtonGroup,
    Checkbox,
    Drawer,
    Heading,
    Input,
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
    Toolbar,
    ToolbarButton,
  } from "flowbite-svelte";
  import {
    ChevronLeftOutline,
    ChevronRightOutline,
    CogSolid,
    DotsVerticalOutline,
    EditOutline,
    ExclamationCircleSolid,
    TrashBinSolid,
  } from "flowbite-svelte-icons";
  import { onMount, type Component } from "svelte";
  import { AddShiftDrawer } from "$lib";
  import type { Shift } from "$lib/model/shift";
  import { deleteShift, fetchShifts } from "$lib/services/shift";
  import { utcToLocal } from "$lib/util/date";

  let shifts: Shift[] = [];
  let loading = true;
  let error = "";

  let searchTerm = "";
  let currentPage = 1;
  const itemsPerPage = 10;
  let totalPages = 1;
  let totalItems = 0;
  let pagesToShow: number[] = [];
  let startRange = 1;
  let endRange = 1;

  let hidden: boolean = true;
  let DrawerComponent: Component = AddShiftDrawer;

  async function handleDeleteShift(id: number) {
    if (!confirm("Are you sure you want to delete this shift?")) return;
    loading = true;
    error = "";
    try {
      const res = await deleteShift(id);
      if (res.success === "success") {
        await loadShifts();
      } else {
        error = "Failed to delete shift";
      }
    } catch (e) {
      error = "Failed to delete shift";
    }
    loading = false;
  }

  function renderPagination() {
    const showPage = 5;
    let startPage = Math.max(1, currentPage - Math.floor(showPage / 2));
    let endPage = Math.min(startPage + showPage - 1, totalPages);
    pagesToShow = Array.from(
      { length: endPage - startPage + 1 },
      (_, i) => startPage + i,
    );
    startRange = shifts.length === 0 ? 0 : (currentPage - 1) * itemsPerPage + 1;
    endRange = Math.min(currentPage * itemsPerPage, totalItems);
  }

  async function loadShifts(page = 1) {
    loading = true;
    error = "";
    try {
      const data = await fetchShifts({
        limit: itemsPerPage,
        offset: (page - 1) * itemsPerPage,
        showOnlyUnassigned: false,
      });
      if (data.success === "success" && data.data) {
        shifts = data.data.shifts;
        currentPage = data.data.meta.current_page;
        totalPages = data.data.meta.total_pages;
        totalItems = data.data.meta.total_items;
        renderPagination();
      } else {
        error = "Failed to fetch shifts";
      }
    } catch (e) {
      error = "Failed to fetch shifts";
    }
    loading = false;
  }

  function loadNextPage() {
    if (currentPage < totalPages) {
      loadShifts(currentPage + 1);
    }
  }

  function loadPreviousPage() {
    if (currentPage > 1) {
      loadShifts(currentPage - 1);
    }
  }

  function goToPage(pageNumber: number) {
    loadShifts(pageNumber);
  }

  function toggle(component: Component) {
    DrawerComponent = component;
    hidden = !hidden;
  }

  onMount(() => loadShifts());
</script>

<main class="relative h-full w-full overflow-y-auto bg-white dark:bg-gray-800">
  <div class="p-4">
    <Breadcrumb class="mb-5">
      <BreadcrumbItem home>Home</BreadcrumbItem>
      <BreadcrumbItem href="/admin/dashboard/shifts">Shifts</BreadcrumbItem>
    </Breadcrumb>
    <Heading
      tag="h1"
      class="text-xl font-semibold text-gray-900 sm:text-2xl dark:text-white"
      >All Shifts</Heading
    >

    <Toolbar embedded class="w-full py-4 text-gray-500 dark:text-gray-300">
      <Input
        placeholder="Search for products"
        class="me-6 w-80 border xl:w-96"
      />
      <ToolbarButton
        color="dark"
        class="m-0 rounded p-1 hover:bg-gray-100 focus:ring-0 dark:hover:bg-gray-700"
      >
        <CogSolid size="lg" />
      </ToolbarButton>
      <ToolbarButton
        color="dark"
        class="m-0 rounded p-1 hover:bg-gray-100 focus:ring-0 dark:hover:bg-gray-700"
      >
        <TrashBinSolid size="lg" />
      </ToolbarButton>
      <ToolbarButton
        color="dark"
        class="m-0 rounded p-1 hover:bg-gray-100 focus:ring-0 dark:hover:bg-gray-700"
      >
        <ExclamationCircleSolid size="lg" />
      </ToolbarButton>
      <ToolbarButton
        color="dark"
        class="m-0 rounded p-1 hover:bg-gray-100 focus:ring-0 dark:hover:bg-gray-700"
      >
        <DotsVerticalOutline size="lg" />
      </ToolbarButton>
      {#snippet end()}
        <div class="space-x-2">
          <Button
            class="whitespace-nowrap"
            onclick={() => toggle(AddShiftDrawer)}>Add new shift</Button
          >
        </div>
      {/snippet}
    </Toolbar>
  </div>
  <Table>
    <TableHead
      class="border-y border-gray-200 bg-gray-100 dark:border-gray-700"
    >
      <TableHeadCell class="w-4 p-4"><Checkbox /></TableHeadCell>
      {#each ["ID", "Date", "Start Time", "End Time", "Role", "Location", "Created At", "Created By", "Actions"] as title}
        <TableHeadCell class="ps-4 font-normal">{title}</TableHeadCell>
      {/each}
    </TableHead>
    <TableBody>
      {#each shifts as shift}
        <TableBodyRow class="text-base">
          <TableBodyCell class="w-4 p-4"><Checkbox /></TableBodyCell>
          <TableBodyCell class="space-x-6 p-4">
            <div class="text-sm font-normal text-gray-500 dark:text-gray-300">
              <div
                class="text-base font-semibold text-gray-900 dark:text-white"
              >
                {shift.id}
              </div>
            </div>
          </TableBodyCell>
          <TableBodyCell class="p-4"
            >{utcToLocal(shift.date, "YYYY-MM-DD")}</TableBodyCell
          >
          <TableBodyCell class="p-4"
            >{utcToLocal(shift.start_time, "HH:mm")}</TableBodyCell
          >
          <TableBodyCell class="p-4"
            >{utcToLocal(shift.end_time, "HH:mm")}</TableBodyCell
          >
          <TableBodyCell class="p-4">{shift.role_id}</TableBodyCell>
          <TableBodyCell class="p-4">{shift.location}</TableBodyCell>
          <TableBodyCell class="p-4"
            >{utcToLocal(shift.created_at, "YYYY-MM-DD HH:mm")}</TableBodyCell
          >
          <TableBodyCell class="p-4">{shift.created_by}</TableBodyCell>
          <TableBodyCell class="space-x-2">
            <Button
              color="green"
              size="sm"
              class="gap-2 px-3"
              onclick={() => toggle(AddShiftDrawer)}
            >
              Assign Worker
            </Button>
            <Button
              size="sm"
              class="gap-2 px-3"
              onclick={() => toggle(AddShiftDrawer)}
            >
              <EditOutline size="sm" /> Update
            </Button>
            <Button
              color="red"
              size="sm"
              class="gap-2 px-3"
              onclick={() => handleDeleteShift(shift.id)}
            >
              <TrashBinSolid size="sm" /> Delete shift
            </Button>
          </TableBodyCell>
        </TableBodyRow>
      {/each}
    </TableBody>
  </Table>
  <div
    class="flex flex-col items-start justify-between space-y-3 p-4 md:flex-row md:items-center md:space-y-0"
    aria-label="Table navigation"
  >
    <span class="text-sm font-normal text-gray-500 dark:text-gray-400">
      Showing
      <span class="font-semibold text-gray-900 dark:text-white"
        >{startRange}-{endRange}</span
      >
      of
      <span class="font-semibold text-gray-900 dark:text-white"
        >{totalItems}</span
      >
    </span>
    <ButtonGroup>
      <Button onclick={loadPreviousPage} disabled={currentPage === 1}
        ><ChevronLeftOutline size="xs" class="m-1.5" /></Button
      >
      {#each pagesToShow as pageNumber}
        <Button
          onclick={() => goToPage(pageNumber)}
          color={pageNumber === currentPage ? "primary" : "gray"}
          >{pageNumber}</Button
        >
      {/each}
      <Button onclick={loadNextPage} disabled={currentPage === totalPages}
        ><ChevronRightOutline size="xs" class="m-1.5" /></Button
      >
    </ButtonGroup>
  </div>
</main>

<Drawer placement="right" bind:hidden>
  <DrawerComponent bind:hidden onSubmit={loadShifts} />
</Drawer>
