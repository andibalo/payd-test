<script lang="ts">
  import {
    Button,
    CloseButton,
    Datepicker,
    Heading,
    Input,
    Label,
  } from "flowbite-svelte";
  import { CloseOutline } from "flowbite-svelte-icons";
  import type { AddShiftDrawerProps } from "./types";
  import { createShift } from "$lib/services/shift";
  import type { CreateShiftReq } from "$lib/request/shift";

  let {
    hidden = $bindable(true),
    title = "Add new shift",
    onSubmit = () => {},
  }: AddShiftDrawerProps = $props();

  let selectedDate = $state<Date | undefined>(undefined);
  let start_time = "";
  let end_time = "";
  let role_id: number = 1;
  let location = "";
  let error = "";

  async function handleCreateShift(event: Event) {
    event.preventDefault();
    error = "";
    if (!selectedDate || !start_time || !end_time || !role_id || !location) {
      error = "Please fill in all fields.";
      return;
    }

    const dateStr = selectedDate.toISOString().split("T")[0];
    const start = new Date(`${dateStr}T${start_time}:00Z`);
    const end = new Date(`${dateStr}T${end_time}:00Z`);
    if (start >= end) {
      error = "Start time must be before end time.";
      return;
    }

    const dateISO = dateStr + "T00:00:00Z";
    const startISO = start.toISOString();
    const endISO = end.toISOString();

    let req: CreateShiftReq = {
      date: dateISO,
      start_time: startISO,
      end_time: endISO,
      role_id,
    };

    if (location) {
      req.location = location;
    }

    const res = await createShift(req);

    if (res.success === "success") {
      hidden = true;

      onSubmit()
    } else {
      error = "Failed to create shift";
    }
  }
</script>

<Heading tag="h5" class="mb-6 text-sm font-semibold uppercase">{title}</Heading>
<CloseButton
  onclick={() => (hidden = true)}
  class="absolute top-2.5 right-2.5 text-gray-400 hover:text-black dark:text-white"
/>

<form on:submit|preventDefault={handleCreateShift}>
  <div class="space-y-4">
    <Label class="space-y-2">
      <span>Date</span>
      <Datepicker bind:value={selectedDate} required />
    </Label>
    <Label class="space-y-2">
      <span>Start Time</span>
      <Input type="time" bind:value={start_time} required />
    </Label>
    <Label class="space-y-2">
      <span>End Time</span>
      <Input type="time" bind:value={end_time} required />
    </Label>
    <Label class="space-y-2">
      <span>Role ID</span>
      <Input type="number" bind:value={role_id} min="1" required />
    </Label>
    <Label class="space-y-2">
      <span>Location</span>
      <Input type="text" bind:value={location} />
    </Label>
    {#if error}
      <div class="text-red-500">{error}</div>
    {/if}
    <div
      class="bottom-0 left-0 flex w-full justify-center space-x-4 pb-4 md:absolute md:px-4"
    >
      <Button type="submit" class="w-full">Add shift</Button>
      <Button
        color="alternative"
        class="w-full"
        type="button"
        onclick={() => (hidden = true)}
      >
        <CloseOutline />
        Cancel
      </Button>
    </div>
  </div>
</form>
