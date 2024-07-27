package dashboard

type PanelType string

const (
	PanelTypeRow        PanelType = "row"
	PanelTypeGauge      PanelType = "gauge"
	PanelTypeStat       PanelType = "stat"
	PanelTypeTimeseries PanelType = "timeseries"
	PanelTypeBargauge   PanelType = "bargauge"
	PanelTypeUnknown    PanelType = "unknown"
)

type (
	Panel struct {
		Datasource struct {
			Type string `json:"type"`
			Uid  string `json:"uid"`
		} `json:"datasource"`
		GridPos GridPos	   `json:"gridPos"`
		ID      int           `json:"id"`
		Panels  []interface{} `json:"panels,omitempty"`
		Targets []struct {
			Datasource struct {
				Type string `json:"type"`
				Uid  string `json:"uid"`
			} `json:"datasource"`
			RefId          string `json:"refId"`
			EditorMode     string `json:"editorMode,omitempty"`
			Exemplar       bool   `json:"exemplar,omitempty"`
			Expr           string `json:"expr,omitempty"`
			Format         string `json:"format,omitempty"`
			Instant        bool   `json:"instant,omitempty"`
			IntervalFactor int    `json:"intervalFactor,omitempty"`
			LegendFormat   string `json:"legendFormat,omitempty"`
			Range          bool   `json:"range,omitempty"`
			Step           int    `json:"step,omitempty"`
			Hide           bool   `json:"hide,omitempty"`
			Interval       string `json:"interval,omitempty"`
		} `json:"targets"`
		Title       string    `json:"title"`
		Type        PanelType `json:"type"`
		Description string    `json:"description,omitempty"`
		FieldConfig struct {
			Defaults struct {
				Color struct {
					Mode string `json:"mode"`
				} `json:"color"`
				Decimals int           `json:"decimals,omitempty"`
				Links    []interface{} `json:"links,omitempty"`
				Mappings []struct {
					Options struct {
						Match  string `json:"match"`
						Result struct {
							Text string `json:"text"`
						} `json:"result"`
					} `json:"options"`
					Type string `json:"type"`
				} `json:"mappings"`
				Max        int `json:"max,omitempty"`
				Min        int `json:"min,omitempty"`
				Thresholds struct {
					Mode  string `json:"mode"`
					Steps []struct {
						Color string `json:"color"`
						Value *int   `json:"value"`
					} `json:"steps"`
				} `json:"thresholds"`
				Unit   string `json:"unit"`
				Custom struct {
					AxisBorderShow   bool   `json:"axisBorderShow"`
					AxisCenteredZero bool   `json:"axisCenteredZero"`
					AxisColorMode    string `json:"axisColorMode"`
					AxisLabel        string `json:"axisLabel"`
					AxisPlacement    string `json:"axisPlacement"`
					BarAlignment     int    `json:"barAlignment"`
					DrawStyle        string `json:"drawStyle"`
					FillOpacity      int    `json:"fillOpacity"`
					GradientMode     string `json:"gradientMode"`
					HideFrom         struct {
						Legend  bool `json:"legend"`
						Tooltip bool `json:"tooltip"`
						Viz     bool `json:"viz"`
					} `json:"hideFrom"`
					InsertNulls       bool   `json:"insertNulls"`
					LineInterpolation string `json:"lineInterpolation"`
					LineWidth         int    `json:"lineWidth"`
					PointSize         int    `json:"pointSize"`
					ScaleDistribution struct {
						Type string `json:"type"`
					} `json:"scaleDistribution"`
					ShowPoints string `json:"showPoints"`
					SpanNulls  bool   `json:"spanNulls"`
					Stacking   struct {
						Group string `json:"group"`
						Mode  string `json:"mode"`
					} `json:"stacking"`
					ThresholdsStyle struct {
						Mode string `json:"mode"`
					} `json:"thresholdsStyle"`
				} `json:"custom,omitempty"`
			} `json:"defaults"`
			Overrides []struct {
				Matcher struct {
					Id      string `json:"id"`
					Options string `json:"options"`
				} `json:"matcher"`
				Properties []struct {
					Id    string      `json:"id"`
					Value interface{} `json:"value"`
				} `json:"properties"`
			} `json:"overrides"`
		} `json:"fieldConfig,omitempty"`
		Options struct {
			DisplayMode   string `json:"displayMode,omitempty"`
			MaxVizHeight  int    `json:"maxVizHeight,omitempty"`
			MinVizHeight  int    `json:"minVizHeight,omitempty"`
			MinVizWidth   int    `json:"minVizWidth,omitempty"`
			NamePlacement string `json:"namePlacement,omitempty"`
			Orientation   string `json:"orientation,omitempty"`
			ReduceOptions struct {
				Calcs  []string `json:"calcs"`
				Fields string   `json:"fields"`
				Values bool     `json:"values"`
			} `json:"reduceOptions,omitempty"`
			ShowUnfilled bool   `json:"showUnfilled,omitempty"`
			Sizing       string `json:"sizing,omitempty"`
			Text         struct {
			} `json:"text,omitempty"`
			ValueMode              string `json:"valueMode,omitempty"`
			ShowThresholdLabels    bool   `json:"showThresholdLabels,omitempty"`
			ShowThresholdMarkers   bool   `json:"showThresholdMarkers,omitempty"`
			ColorMode              string `json:"colorMode,omitempty"`
			GraphMode              string `json:"graphMode,omitempty"`
			JustifyMode            string `json:"justifyMode,omitempty"`
			PercentChangeColorMode string `json:"percentChangeColorMode,omitempty"`
			ShowPercentChange      bool   `json:"showPercentChange,omitempty"`
			TextMode               string `json:"textMode,omitempty"`
			WideLayout             bool   `json:"wideLayout,omitempty"`
			Legend                 struct {
				Calcs       []string `json:"calcs"`
				DisplayMode string   `json:"displayMode"`
				Placement   string   `json:"placement"`
				ShowLegend  bool     `json:"showLegend"`
				Width       int      `json:"width,omitempty"`
			} `json:"legend,omitempty"`
			Tooltip struct {
				Mode string `json:"mode"`
				Sort string `json:"sort"`
			} `json:"tooltip,omitempty"`
		} `json:"options,omitempty"`
		PluginVersion    string `json:"pluginVersion,omitempty"`
		HideTimeOverride bool   `json:"hideTimeOverride,omitempty"`
		MaxDataPoints    int    `json:"maxDataPoints,omitempty"`
	}

	Dashboard struct {
		Annotations struct {
			List []struct {
				HashKey    string `json:"$$hashKey"`
				BuiltIn    int    `json:"builtIn"`
				Datasource struct {
					Type string `json:"type"`
					Uid  string `json:"uid"`
				} `json:"datasource"`
				Enable    bool   `json:"enable"`
				Hide      bool   `json:"hide"`
				IconColor string `json:"iconColor"`
				Name      string `json:"name"`
				Target    struct {
					Limit    int           `json:"limit"`
					MatchAny bool          `json:"matchAny"`
					Tags     []interface{} `json:"tags"`
					Type     string        `json:"type"`
				} `json:"target"`
				Type string `json:"type"`
			} `json:"list"`
		} `json:"annotations"`
		Editable             bool     `json:"editable"`
		FiscalYearStartMonth int      `json:"fiscalYearStartMonth"`
		GnetId               int      `json:"gnetId"`
		GraphTooltip         int      `json:"graphTooltip"`
		Id                   int      `json:"id"`
		Panels               []*Panel `json:"panels"`
		Refresh              string   `json:"refresh"`
		Revision             int      `json:"revision"`
		SchemaVersion        int      `json:"schemaVersion"`
		Tags                 []string `json:"tags"`
		Templating           struct {
			List []struct {
				Current struct {
					Selected bool   `json:"selected"`
					Text     string `json:"text"`
					Value    string `json:"value"`
				} `json:"current"`
				Hide       int    `json:"hide"`
				IncludeAll bool   `json:"includeAll"`
				Label      string `json:"label,omitempty"`
				Multi      bool   `json:"multi"`
				Name       string `json:"name"`
				Options    []struct {
					Selected bool   `json:"selected"`
					Text     string `json:"text"`
					Value    string `json:"value"`
				} `json:"options"`
				Query       interface{} `json:"query"`
				QueryValue  string      `json:"queryValue,omitempty"`
				Refresh     int         `json:"refresh,omitempty"`
				Regex       string      `json:"regex,omitempty"`
				SkipUrlSync bool        `json:"skipUrlSync"`
				Type        string      `json:"type"`
				Datasource  struct {
					Type string `json:"type"`
					Uid  string `json:"uid"`
				} `json:"datasource,omitempty"`
				Definition     string `json:"definition,omitempty"`
				Sort           int    `json:"sort,omitempty"`
				TagValuesQuery string `json:"tagValuesQuery,omitempty"`
				TagsQuery      string `json:"tagsQuery,omitempty"`
				UseTags        bool   `json:"useTags,omitempty"`
			} `json:"list"`
		} `json:"templating"`
		Time struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"time"`
		Timepicker struct {
			RefreshIntervals []string `json:"refresh_intervals"`
			TimeOptions      []string `json:"time_options"`
		} `json:"timepicker"`
		Timezone  string `json:"timezone"`
		Title     string `json:"title"`
		Uid       string `json:"uid"`
		Version   int    `json:"version"`
		WeekStart string `json:"weekStart"`
	}
)
