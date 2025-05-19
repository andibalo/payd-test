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
    ExclamationCircleSolid,
    TrashBinSolid,
  } from "flowbite-svelte-icons";
  import { onMount, type Component } from "svelte";
  import { AddShiftDrawer } from "$lib";
  import {
    approveShiftRequest,
    deleteShift,
    fetchShiftRequests,
    rejectShiftRequest,
  } from "$lib/services/shift";
  import { utcToLocal } from "$lib/util/date";
  import type { ShiftRequest } from "$lib/response/shift";

  let shiftRequests: ShiftRequest[] = [];
  let loading = true;
  let error = "";

  let searchTerm = "";
  let status = "";
  let currentPage = 1;
  const itemsPerPage = 10;
  let totalPages = 1;
  let totalItems = 0;
  let pagesToShow: number[] = [];
  let startRange = 1;
  let endRange = 1;

  let hidden: boolean = true;
  let DrawerComponent: Component = AddShiftDrawer;

  function renderPagination() {
    const showPage = 5;
    let startPage = Math.max(1, currentPage - Math.floor(showPage / 2));
    let endPage = Math.min(startPage + showPage - 1, totalPages);
    pagesToShow = Array.from(
      { length: endPage - startPage + 1 },
      (_, i) => startPage + i,
    );
    startRange =
      shiftRequests.length === 0 ? 0 : (currentPage - 1) * itemsPerPage + 1;
    endRange = Math.min(currentPage * itemsPerPage, totalItems);
  }

  async function loadShiftRequests(page = 1) {
    loading = true;
    error = "";
    try {
      const data = await fetchShiftRequests({
        limit: itemsPerPage,
        offset: (page - 1) * itemsPerPage,
        status: status,
      });
      if (data.success === "success" && data.data) {
        shiftRequests = data.data.request_shifts;
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

  async function handleApproveShiftRequest(id: number) {
    loading = true;
    error = "";
    try {
      const res = await approveShiftRequest(id);
      if (res.success === "success") {
        await loadShiftRequests(currentPage);
      } else {
        error = "Failed to approve shift request";
      }
    } catch (e) {
      error = "Failed to approve shift request";
    }
    loading = false;
  }

  async function handleRejectShiftRequest(id: number) {
    const reason = prompt("Enter rejection reason:");
    if (!reason) return;
    loading = true;
    error = "";
    try {
      const res = await rejectShiftRequest(id, reason);
      if (res.success === "success") {
        await loadShiftRequests(currentPage);
      } else {
        error = "Failed to reject shift request";
      }
    } catch (e) {
      error = "Failed to reject shift request";
    }
    loading = false;
  }

  function loadNextPage() {
    if (currentPage < totalPages) {
      loadShiftRequests(currentPage + 1);
    }
  }

  function loadPreviousPage() {
    if (currentPage > 1) {
      loadShiftRequests(currentPage - 1);
    }
  }

  function goToPage(pageNumber: number) {
    loadShiftRequests(pageNumber);
  }

  function toggle(component: Component) {
    DrawerComponent = component;
    hidden = !hidden;
  }

  onMount(() => loadShiftRequests());
</script>

<main class="relative h-full w-full overflow-y-auto bg-white dark:bg-gray-800">
  <div class="p-4">
    <Breadcrumb class="mb-5">
      <BreadcrumbItem home>Home</BreadcrumbItem>
      <BreadcrumbItem href="/admin/dashboard/shift-requests"
        >Shift Reqeusts</BreadcrumbItem
      >
    </Breadcrumb>
    <Heading
      tag="h1"
      class="text-xl font-semibold text-gray-900 sm:text-2xl dark:text-white"
      >All Shift Requests</Heading
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
    </Toolbar>
  </div>
  <Table>
    <TableHead
      class="border-y border-gray-200 bg-gray-100 dark:border-gray-700"
    >
      <TableHeadCell class="w-4 p-4"><Checkbox /></TableHeadCell>
      {#each ["ID", "Shift Date", "Start Time", "End Time", "Role", "Status", "Requested By", "Admin Actor", "Rejection Reason", "Created At", "Actions"] as title}
        <TableHeadCell class="ps-4 font-normal">{title}</TableHeadCell>
      {/each}
    </TableHead>
    <TableBody>
      {#each shiftRequests as shiftRequest}
        <TableBodyRow class="text-base">
          <TableBodyCell class="w-4 p-4"><Checkbox /></TableBodyCell>
          <TableBodyCell class="space-x-6 p-4">
            <div class="text-sm font-normal text-gray-500 dark:text-gray-300">
              <div
                class="text-base font-semibold text-gray-900 dark:text-white"
              >
                {shiftRequest.id}
              </div>
            </div>
          </TableBodyCell>
          <TableBodyCell class="p-4"
            >{utcToLocal(shiftRequest.shift_date, "YYYY-MM-DD")}</TableBodyCell
          >
          <TableBodyCell class="p-4"
            >{utcToLocal(shiftRequest.shift_start_time, "HH:mm")}</TableBodyCell
          >
          <TableBodyCell class="p-4"
            >{utcToLocal(shiftRequest.shift_end_time, "HH:mm")}</TableBodyCell
          >
          <TableBodyCell class="p-4"
            >{shiftRequest.shift_role_name}</TableBodyCell
          >
          <TableBodyCell class="p-4">{shiftRequest.status}</TableBodyCell>
          <TableBodyCell class="p-4">{shiftRequest.requested_by}</TableBodyCell>
          <TableBodyCell class="p-4">{shiftRequest.admin_actor}</TableBodyCell>
          <TableBodyCell class="p-4"
            >{shiftRequest.rejection_reason}</TableBodyCell
          >
          <TableBodyCell class="p-4"
            >{utcToLocal(
              shiftRequest.created_at,
              "YYYY-MM-DD HH:mm",
            )}</TableBodyCell
          >
          <TableBodyCell class="space-x-2">
            <Button
              color="green"
              size="sm"
              class="gap-2 px-3"
              onclick={() => handleApproveShiftRequest(shiftRequest.id)}
            >
              Approve
            </Button>
            <Button
              color="red"
              size="sm"
              class="gap-2 px-3"
              onclick={() => handleRejectShiftRequest(shiftRequest.id)}
            >
              <TrashBinSolid size="sm" /> Reject
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
  <DrawerComponent bind:hidden onSubmit={loadShiftRequests} />
</Drawer>
