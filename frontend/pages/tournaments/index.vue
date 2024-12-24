<template>
  <div>
    <Card>
      <CardHeader>
        <CardTitle>Tournaments</CardTitle>
      </CardHeader>
    </Card>
    <Card>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Tier</TableHead>
              <TableHead>Location</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="tournament in tournaments.data" :key="tournament.ExternalID">
              <TableCell class="font-medium"> {{ tournament.Name }} </TableCell>
              <TableCell> {{ tournament.Tier }} </TableCell>
              <TableCell> {{tournament.Location}} </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import type { ListResult, TournamentRow } from '@/types/api'

const { data } = await useFetch('/api/tournaments', {server: true})
const tournaments = data.value as ListResult<TournamentRow>
const searchTerm = ref("")
</script>
