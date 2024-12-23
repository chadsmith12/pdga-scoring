<template>
  <div class="container mx-auto p-4">
    <Card>
      <CardHeader>
        <CardTitle>Player Performance</CardTitle>
      </CardHeader>
      <CardContent>
        <div class="flex items-center space-x-2">
          <Input v-model="searchTerm" type="text" placeholder="Search Players..." class="max-w-sm"></Input>
        </div>
      </CardContent>
    </Card>
    <Card>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Division</TableHead>
              <TableHead>Avg. Points</TableHead>
              <TableHead>Trend</TableHead>
              <TableHead>Wins</TableHead>
              <TableHead>Podiums</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="player in players" :key="player.id">
              <TableCell class="font-medium">
                {{ player.name }}
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
interface Player {
  id: string
  name: string
  division: 'MPO' | 'FPO'
  averagePoints: number
  trend: 'up' | 'down' | 'stable'
  wins: number
  podiums: number
}

const searchTerm = ref("")
const players = ref<Player[]>([
  { id: '1', name: 'Paul McBeth', division: 'MPO', averagePoints: 95.5, trend: 'up', wins: 3, podiums: 5 },
  { id: '2', name: 'Ricky Wysocki', division: 'MPO', averagePoints: 93.2, trend: 'stable', wins: 2, podiums: 4 },
  { id: '3', name: 'Eagle McMahon', division: 'MPO', averagePoints: 91.8, trend: 'down', wins: 1, podiums: 3 },
  { id: '4', name: 'Paige Pierce', division: 'FPO', averagePoints: 94.7, trend: 'up', wins: 4, podiums: 6 },
  { id: '5', name: 'Catrina Allen', division: 'FPO', averagePoints: 92.1, trend: 'stable', wins: 2, podiums: 5 },
  { id: '6', name: 'Kristin Tattar', division: 'FPO', averagePoints: 93.5, trend: 'up', wins: 3, podiums: 5 },
  { id: '7', name: 'Calvin Heimburg', division: 'MPO', averagePoints: 90.9, trend: 'up', wins: 1, podiums: 4 },
  { id: '8', name: 'James Conrad', division: 'MPO', averagePoints: 89.5, trend: 'down', wins: 1, podiums: 2 },
  { id: '9', name: 'Hailey King', division: 'FPO', averagePoints: 90.8, trend: 'up', wins: 1, podiums: 3 },
  { id: '10', name: 'Simon Lizotte', division: 'MPO', averagePoints: 91.2, trend: 'stable', wins: 2, podiums: 4 },
])
</script>
