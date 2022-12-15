package test_client2.models;

import com.squareup.moshi.Json;
import com.squareup.moshi.Moshi;
import com.squareup.moshi.Types;
import java.math.BigDecimal;
import java.time.*;
import java.util.*;
import java.io.*;


public enum Choice {
	@Json(name = "FIRST_CHOICE") FIRST_CHOICE,
	@Json(name = "SECOND_CHOICE") SECOND_CHOICE,
	@Json(name = "THIRD_CHOICE") THIRD_CHOICE,
}