package schemas

import "gitlab.sandmanbb.com/perfil-digital-agro/agro-ws/build/scicrop-scraper/models/schema"

var (
	Owner = &schema.Schema{ Uri: "scicrop_owner", Fields: []schema.Field{
			{Name: "name", Type: "short_string", Description: "Nome do dono da conta em scicrop.com. Entidade personEntityLst da api."},
			{Name: "phone", Type: "short_string", Description: "Telefone fornecido pelo dono da conta à scicrop. Entidade personEntityLst da api."},
			{Name: "address", Type: "short_string", Description: "Endereço fornecido pelo dono da conta à scicrop. Entidade personEntityLst da api."},
		},
	}
	Thing = &schema.Schema{ Uri: "scicrop_station", Fields: []schema.Field{
			{Name: "name", Type: "short_string", Description: "Nome da estação metereológica em scicrop.com (Mais próximo de um código). Obtido da entidade stationLst da api."},
			{Name: "ref_name", Type: "short_string", Description: "Nome de referência da estação na scicrop (Nome amigável). Obtido da entidade stationLst da api."},
			{Name: "station_id", Type: "short_string", Description: "Id da estação na scicrop. Obtido da entidade stationLst da api."},
		},
	}
	Node = &schema.Schema{ Uri: "scicrop_station_node", Fields: []schema.Field{
			{Name: "prec", Type: "decimal", Description: "Valor diário de precipitação. Obtido do end-point getScicropStationDataByDayOnStationDataT."},
			{Name: "humidity", Type: "decimal", Description: "Valor diário de humidade. Obtido do end-point getScicropStationDataByDayOnStationDataT."},
			{Name: "windDir", Type: "decimal", Description: "Valor diário de direção do vento. Obtido do end-point getScicropStationDataByDayOnStationDataT."},
			{Name: "max_temperature", Type: "decimal", Description: "Máxima diária de temperatura em Graus Celsius. Obtido do end-point getScicropStationDataByDayOnStationDataT."},
			{Name: "min_temperature", Type: "decimal", Description: "Mínima diária de temperatura em Graus Celsius. Obtido do end-point getScicropStationDataByDayOnStationDataT."},
			{Name: "avg_temperature", Type: "decimal", Description: "Média diária de temperatura em Graus Celsius. Obtido do end-point getScicropStationDataByDayOnStationDataT."},
			{Name: "radiation", Type: "decimal", Description: "Valor diário de radiação solar. Obtido do end-point getScicropStationDataByDayOnStationDataT."},
			{Name: "uvIdx", Type: "decimal", Description: "Valor diário de radição ultra violeta. Obtido do end-point getScicropStationDataByDayOnStationDataT."},
			{Name: "readings", Type: "decimal", Description: "Número de leituras da estação no dia. Obtido do end-point getScicropStationDataByDayOnStationDataT."},
		},
	}
	Temp24hs = &schema.Schema{ Uri: "wheater_daily_temperature", Node: "temp24hs",  Fields: []schema.Field{
			{Name: "max_temperature", Type: "decimal", Description: "Máxima diária de temperatura em Graus Celsius."},
			{Name: "min_temperature", Type: "decimal", Description: "Mínima diária de temperatura em Graus Celsius."},
			{Name: "avg_temperature", Type: "decimal", Description: "Média diária de temperatura em Graus Celsius."},
		},
	}
	Wind24hs = &schema.Schema{ Uri: "wheater_daily_wind", Node: "wind24hs", Fields: []schema.Field{
			{Name: "wind_speed", Type: "decimal", Description: "Valor diário de velocidade do vento.", From: "windSpeed"},
			{Name: "wind_gust", Type: "decimal", Description: "Pico diário de intesidade de vento.", From: "windGust"},
			{Name: "wind_direction", Type: "decimal", Description: "Pico diário de intesidade de vento.", From: "windDir"},
		},
	}
	Radiation24hs = &schema.Schema{ Uri: "wheater_daily_radiation", Node: "radiation24hs", Fields: []schema.Field{
			{Name: "radiation", Type: "decimal", Description: "Valor diário de radiação solar."},
			{Name: "uv_index", Type: "decimal", Description: "Valor diário de radição ultra violeta.", From: "uvIdx"},
		},
	}
	Humidity24hs = &schema.Schema{ Uri: "wheater_daily_humidity", Node: "humidity24hs", Fields: []schema.Field{
			{Name: "precipitation", Type: "decimal", Description: "Valor diário de precipitação.", From: "prec"},
			{Name: "humidity", Type: "decimal", Description: "Valor diário de humidade."},
		},
	}
)
