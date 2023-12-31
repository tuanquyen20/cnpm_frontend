<template>
  <v-container fluid>
    <div>
      <v-dialog v-model="createDialog" width="500">
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
            Create User
          </v-card-title>

          <v-card-actions>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  v-model="createForm.avatar"
                  label="Avatar link"
                ></v-text-field>
              </v-col>

              <v-col cols="12">
                <v-text-field
                  v-model="createForm.name"
                  label="Name"
                ></v-text-field>
              </v-col>

              <v-col cols="12">
                <v-text-field
                  v-model="createForm.phone_number"
                  label="Phone Number"
                ></v-text-field>
              </v-col>

              <v-col cols="12">
                <v-text-field
                  v-model="createForm.email"
                  label="Email"
                ></v-text-field>
              </v-col>

              <v-col cols="12">
                <v-text-field
                  v-model="createForm.password"
                  label="Password"
                ></v-text-field>
              </v-col>

              <v-col cols="12">
                <v-select
                  v-model="createForm.type"
                  :items="['admin', 'customer']"
                  label="Type"
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

      <v-text-field v-model="search" label="Search" append-icon="mdi-magnify"></v-text-field>

      <v-data-table
        :headers="headers"
        :items="items"
        :page.sync="page"
        :items-per-page="itemsPerPage"
        hide-default-footer
        class="elevation-1"
        :loading="loading"
      >
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
          Update User
        </v-card-title>

        <v-card-actions>
          <v-row>
            <v-col cols="12">
              <v-text-field
                v-model="updateForm.avatar"
                label="Avatar Link"
              ></v-text-field>
            </v-col>

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

            <v-col cols="12">
              <v-text-field
                v-model="updateForm.password"
                label="Password"
              ></v-text-field>
            </v-col>

            <v-col cols="12">
              <v-select
                v-model="updateForm.type"
                :items="['admin', 'customer']"
                label="Type"
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
      search: "",

      editDialog: false,
      createDialog: false,
      loading: false,
      page: 1,
      pageCount: 0,
      itemsPerPage: 10,
      items: [],
      headers: [
        { text: "ID", value: "id" },
        { text: "Email", value: "email" },
        { text: "Name", value: "name" },
        { text: "Phone Number", value: "phone_number" },
        { text: "Type", value: "type" },
        { text: "Actions", value: "actions", sortable: false },
      ],

      updateForm: {
        id: null,
        name: "",
        email: "",
        password: "",
        phone_number: "",
        avatar: "",
      },

      createForm: {
        name: "",
        email: "",
        password: "",
        phone_number: "",
        avatar: "",
      },
    };
  },

  watch: {
    search(newValue) {
      this.loadItems()
    }
  },

  methods: {
    loadItems() {
      this.loading = true;
      setTimeout(() => {
        this.loading = false;
        this.$axios
          .get(`/users?page=${this.page}&limit=${this.itemsPerPage}&search=${this.search}`)
          .then((res) => {
            this.items = res.data.items || [];
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
        .delete(`/users/${item.id}`)
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
        email: item.email,
        password: item.password,
        phone_number: item.phone_number,
        type: item.type,
        avatar: item.avatar
      };
    },

    updateItem() {
      this.$axios
        .put(`/users/${this.updateForm.id}`, {
          name: this.updateForm.name,
          phone_number: this.updateForm.phone_number,
          email: this.updateForm.email,
          password: this.updateForm.password,
          type: this.updateForm.type,
          avatar: this.updateForm.avatar
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
        .post(`/users`, {
          name: this.createForm.name,
          phone_number: this.createForm.phone_number,
          email: this.createForm.email,
          password: this.createForm.password,
          type: this.createForm.type,
          avatar: this.createForm.avatar,
        })

        .then((res) => {
          this.createDialog = false;
          this.loadItems();
        })
        .catch((err) => {
          console.error(err);
        });
    },
  },

  created() {
    this.loadItems();
  },
};
</script>
