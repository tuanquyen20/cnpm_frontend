<template>
  <v-container fluid>
    <div>
      <v-dialog v-model="createDialog" width="900">
        <template v-slot:activator="{ on, attrs }">
          <v-btn
            dark
            v-bind="attrs"
            v-on="on"
            small
            color="#252525"
            class="mb-4"
          >
            <v-icon dark> mdi-plus </v-icon>
            Create
          </v-btn>
        </template>
        <v-card>
          <v-card-title class="text-h5 grey lighten-2">
            Create Room
          </v-card-title>

          <v-card-actions>
            <v-row>
              <v-col cols="12">
                <v-select
                  v-model="createForm.select"
                  :hint="`ID: ${createForm.select.id}`"
                  :items="roomTypesList"
                  item-text="name"
                  item-value="id"
                  label="Room Type"
                  persistent-hint
                  return-object
                  single-line
                ></v-select>
              </v-col>
            </v-row>
          </v-card-actions>

          <v-divider></v-divider>

          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="#252525" class="white--text" @click="createItem">
              Create
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>

      <v-data-table
        :headers="headers"
        :items="items"
        :page.sync="page"
        :items-per-page="itemsPerPage"
        hide-default-footer
        class="elevation-1"
        :loading="loading"
      >
        <template v-slot:[`item.room_type_id`]="{ item }">
          <span>
            {{ renderRoomType(item.room_type_id) }}
          </span>
        </template>

        <template v-slot:[`item.actions`]="{ item }">
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

    <v-dialog v-model="editDialog" width="900">
      <v-card>
        <v-card-title class="text-h5 grey lighten-2">
          Update Room Type
        </v-card-title>

        <v-card-actions>
          <v-row>
            <v-col cols="12">
              <v-select
                v-model="updateForm.select"
                :hint="`ID: ${updateForm.select.id}`"
                :items="roomTypesList"
                item-text="name"
                item-value="id"
                label="Room Type"
                persistent-hint
                return-object
                single-line
              ></v-select>
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
import numeral from "numeral";

export default {
  data() {
    return {
      roomTypesList: [],
      editDialog: false,
      createDialog: false,
      loading: false,
      page: 1,
      pageCount: 0,
      itemsPerPage: 10,
      items: [],
      headers: [
        { text: "ID", value: "id" },
        { text: "Room Type", value: "room_type_id" },

        { text: "Actions", value: "actions", sortable: false },
      ],

      updateForm: {
        id: null,
        select: {
          id: "",
          name: "",
        },
      },

      createForm: {
        select: {
          id: "",
          name: "",
        },
      },
    };
  },

  methods: {
    loadItems() {
      this.loading = true;
      setTimeout(() => {
        this.loading = false;
        this.$axios
          .get(`/rooms?page=${this.page}&limit=${this.itemsPerPage}`)
          .then((res) => {
            this.items = res.data.items;
            this.pageCount = res.data.page_count;
          })
          .catch((err) => {
            console.error(err);
          });
      }, 1000);
    },

    formattedPrice(price) {
      return numeral(price).format("0,0");
    },

    deleteItem(item) {
      this.$axios
        .delete(`/rooms/${item.id}`)
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
        select: {
          id: item.room_type_id,
          name: this.renderRoomType(item.room_type_id),
        },
      };
    },

    updateItem() {
      this.$axios
        .put(`/rooms/${this.updateForm.id}`, {
          room_type_id: this.updateForm.select.id,
        })

        .then((res) => {
          this.editDialog = false;
          this.loadItems();
        })
        .catch((err) => {
          console.error(err);
        });
    },

    createItem() {
      this.$axios
        .post(`/rooms`, {
          room_type_id: this.createForm.select.id,
        })

        .then((res) => {
          this.createDialog = false;
          this.createForm = {
            select: {
              id: "",
              name: "",
            },
          };
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
