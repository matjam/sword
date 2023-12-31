package assets

type tilePosition struct {
	name string
	x    int
	y    int
}

var dungeonTiles = []tilePosition{
	{"blue_wall_corner_top_left", 0, 0},
	{"blue_wall_cross", 1, 0},
	{"blue_wall_horizontal_1", 2, 0},
	{"blue_wall_horizontal_2", 4, 2},
	{"blue_wall_horizontal_2", 4, 2},
	{"blue_wall_horizontal_end_left", 1, 1},
	{"blue_wall_horizontal_end_right", 2, 1},
	{"blue_wall_horizontal_rubble_left", 4, 3},
	{"blue_wall_horizontal_rubble_right", 5, 3},
	{"blue_wall_corner_top_right", 3, 0},
	{"blue_wall_junction_right", 0, 1},
	{"blue_wall_junction_left", 3, 1},
	{"blue_wall_vertical_1", 0, 2},
	{"blue_wall_vertical_2", 3, 2},
	{"blue_wall_vertical_3", 4, 0},
	{"blue_wall_vertical_end_up", 1, 2},
	{"blue_wall_vertical_end_down", 2, 3},
	{"blue_wall_junction_down", 2, 2},
	{"blue_wall_corner_bottom_left", 0, 3},
	{"blue_wall_junction_up", 0, 3},
	{"blue_wall_pillar", 4, 1},
	{"blue_rubble_1", 5, 3},
	{"blue_rubble_2", 6, 2},
	{"blue_rubble_3", 7, 2},
	{"blue_tile_1", 5, 0},
	{"blue_tile_2", 6, 0},
	{"blue_tile_3", 7, 0},
	{"blue_tile_4", 5, 1},
	{"blue_tile_5", 6, 1},
	{"blue_tile_6", 7, 1},
	{"gray_wall_corner_top_left", 8, 0},
	{"gray_wall_cross", 9, 0},
	{"gray_wall_horizontal_1", 10, 0},
	{"gray_wall_horizontal_2", 12, 2},
	{"gray_wall_horizontal_2", 12, 2},
	{"gray_wall_horizontal_end_left", 9, 1},
	{"gray_wall_horizontal_end_right", 10, 1},
	{"gray_wall_horizontal_rubble_left", 12, 3},
	{"gray_wall_horizontal_rubble_right", 13, 3},
	{"gray_wall_corner_top_right", 11, 0},
	{"gray_wall_junction_right", 8, 1},
	{"gray_wall_junction_left", 11, 1},
	{"gray_wall_vertical_1", 8, 2},
	{"gray_wall_vertical_2", 11, 2},
	{"gray_wall_vertical_3", 12, 0},
	{"gray_wall_vertical_end_up", 9, 2},
	{"gray_wall_vertical_end_down", 10, 3},
	{"gray_wall_junction_down", 10, 2},
	{"gray_wall_corner_bottom_left", 8, 3},
	{"gray_wall_junction_up", 8, 3},
	{"gray_wall_pillar", 12, 1},
	{"gray_rubble_1", 13, 3},
	{"gray_rubble_2", 14, 2},
	{"gray_rubble_3", 15, 2},
	{"gray_tile_1", 13, 0},
	{"gray_tile_2", 14, 0},
	{"gray_tile_3", 15, 0},
	{"gray_tile_4", 13, 1},
	{"gray_tile_5", 14, 1},
	{"gray_tile_6", 15, 1},
}
