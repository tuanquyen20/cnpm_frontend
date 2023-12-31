<template>
  <v-container fluid>
    <div>
      <v-text-field
        v-model="search"
        label="Search"
        append-icon="mdi-magnify"
      ></v-text-field>
      <v-data-table
        :headers="headers"
        :items="items"
        :page.sync="page"
        :items-per-page="itemsPerPage"
        hide-default-footer
        class="elevation-1"
        :loading="loading"
      >
        <template v-slot:[`item.time`]="{ item }">
          <span>
            {{ item.time_start + " - " + item.time_end }}
          </span>
        </template>

        <template v-slot:[`item.room_type`]="{ item }">
          <span>
            {{ renderRoomType(item.room_type_id) }}
          </span>
        </template>

        <template v-slot:[`item.phone_number`]="{ item }">
          <strong>
            {{ item.phone_number }}
          </strong>
        </template>

        <template v-slot:[`item.status`]="{ item }">
          <v-chip
            :color="item.status === 'unconfirmed' ? 'red' : 'green'"
            text-color="white"
          >
            {{ item.status === "unconfirmed" ? "Unconfirmed" : "Confirmed" }}
          </v-chip>
        </template>

        <template v-slot:[`item.actions`]="{ item }">
          <v-icon
            v-if="item.status === 'unconfirmed'"
            small
            color="green"
            class="mr-2"
            @click="updateStatus(item)"
          >
            mdi-check-circle
          </v-icon>

          <v-icon
            v-if="item.status === 'confirmed'"
            small
            color="red"
            class="mr-2"
            @click="updateStatus(item)"
          >
            mdi-close-box
          </v-icon>

          <v-icon small class="mr-2" @click="showUpdateForm(item)">
            mdi-pencil
          </v-icon>

          <v-icon small @click="deleteItem(item)"> mdi-delete </v-icon>
        </template>
      </v-data-table>

      <v-pagination
        v-model="page"
        :length="pageCount"
        @next="loadItems"
        @previous="loadItems"
        @input="loadItems"
      ></v-pagination>
    </div>

    <v-dialog v-model="editDialog" width="500">
      <v-card>
        <v-card-title class="text-h5 grey lighten-2">
          Update Order
        </v-card-title>

        <v-card-actions>
          <v-row>
            <v-col cols="12">
              <v-text-field
                v-model="updateForm.name"
                label="Name"
              ></v-text-field>
            </v-col>

            <v-col cols="12">
              <v-text-field
                v-model="updateForm.phone_number"
                label="Phone Number"
              ></v-text-field>
            </v-col>

            <v-col cols="12">
              <v-text-field
                v-model="updateForm.email"
                label="Email"
              ></v-text-field>
            </v-col>
          </v-row>
        </v-card-actions>

        <v-divider></v-divider>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="#252525" class="white--text" @click="updateItem">
            Update
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script>
export default {
  data() {
    return {
      editDialog: false,
      search: "",
      roomTypesList: [],
      loading: false,
      page: 1,
      pageCount: 0,
      itemsPerPage: 10,
      items: [],
      headers: [
        { text: "ID", value: "id" },
        { text: "Name", value: "name" },
        { text: "Phone Number", value: "phone_number" },
        { text: "Email", value: "email" },
        { text: "Time", value: "time" },
        { text: "Room ID", value: "room_id" },
        { text: "Status", value: "status" },
        { text: "Room Type", value: "room_type" },
        { text: "Actions", value: "actions", sortable: false },
      ],

      updateForm: {
        id: null,
        name: "",
        email: "",
        phone_number: "",
      },
    };
  },

  watch: {
    search(newValue) {
      this.loadItems();
    },
  },

  methods: {
    loadItems() {
      this.loading = true;
      setTimeout(() => {
        this.loading = false;
        this.$axios
          .get(
            `/orders?page=${this.page}&limit=${this.itemsPerPage}&search=${this.search}`
          )
          .then((res) => {
            this.items = res.data.items || [];
            this.pageCount = res.data.page_count;
          })
          .catch((err) => {
            console.error(err);
          });
      }, 1000);
    },

    deleteItem(item) {
      this.$axios
        .delete(`/orders/${item.id}`)
        .then((res) => {
          this.loadItems();
        })
        .catch((err) => {
          console.error(err);
        });
    },

    showUpdateForm(item) {
      this.editDialog = true;
      this.updateForm = {
        id: item.id,
        name: item.name,
        phone_number: item.phone_number,
        email: item.email,
        time_start: item.time_start,
        time_end: item.time_end,
        room_id: item.room_id,
        status: item.status,
        room_type_id: item.room_type_id,
      };
    },

    updateStatus(item) {
      this.$axios
        .put(`/orders/${item.id}`, {
          id: item.id,
          name: item.name,
          phone_number: item.phone_number,
          email: item.email,
          time_start: item.time_start,
          time_end: item.time_end,
          room_id: item.room_id,
          status: item.status === "unconfirmed" ? "confirmed" : "unconfirmed",
          room_type_id: item.room_type_id,
        })

        .then((res) => {
          this.loadItems();
        });
    },

    updateItem() {
      this.$axios
        .put(`/orders/${this.updateForm.id}`, {
          id: this.updateForm.id,
          name: this.updateForm.name,
          phone_number: this.updateForm.phone_number,
          email: this.updateForm.email,
          time_start: this.updateForm.time_start,
          time_end: this.updateForm.time_end,
          room_id: this.updateForm.room_id,
          status: this.updateForm.status,
          room_type_id: this.updateForm.room_type_id,
        })

        .then((res) => {
          this.editDialog = false;
          this.loadItems();
        })
        .catch((err) => {
          console.error(err);
        });
    },

    async getRoomTypesList() {
      await this.$axios
        .get("/room_types/all")
        .then((res) => {
          this.roomTypesList = res.data;
        })
        .catch((err) => {
          console.error(err);
        });
    },

    renderRoomType(id) {
      const x = this.roomTypesList.find((i) => i.id === id);
      return x.name;
    },
  },

  async created() {
    await this.getRoomTypesList();
    this.loadItems();
  },
};
</script>
